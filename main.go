package main

import (
	"embed"
	"flag"
	"fmt"
	"gostrecka/services/database/sqlite"
	"gostrecka/services/discord/commands"
	"gostrecka/services/env"
	"gostrecka/services/transactions"
	"gostrecka/utils"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"
	"github.com/sarulabs/di/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/zekrotja/ken"
)

var nFlag = flag.Bool("v", false, "Version")

//go:embed all:frontend/dist
var Assets embed.FS

func main() {
	flag.Parse()
	if *nFlag {
		fmt.Println("Version: 0.0.1")
		return
	}

	runtime.LockOSThread()

	builder, err := di.NewEnhancedBuilder()
	if err != nil {
		panic(err)
	}

	builder.Add(&di.Def{

		Name: "logger",
		Build: func(ctn di.Container) (interface{}, error) {
			return NewLogger(), nil
		},
	})

	builder.Add(&di.Def{
		Name: "config",
		Build: func(ctn di.Container) (interface{}, error) {
			logger := ctn.Get("logger").(*slog.Logger)
			configStore, err := env.NewConfigStore(logger.With("service", "CONFIG"))
			if err != nil {
				return nil, err
			}
			return configStore.Config()
		},
		Close: func(obj interface{}) error {
			// If your Config needs any cleanup, do it here
			return nil
		},
	})

	builder.Add(&di.Def{
		Name: "database",
		Build: func(ctn di.Container) (interface{}, error) {
			db := sqlite.New(ctn)
			err := db.Connect()

			if err != nil {
				return nil, err
			}

			return db, nil
		},
		Close: func(obj interface{}) error {
			// If your Database needs any cleanup, do it here
			return nil
		},
	})

	builder.Add(&di.Def{Name: "discord_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return NewDiscordService(ctn), nil
		},
		Close: func(obj interface{}) error {
			service := obj.(*DiscordService)
			return service.session.Close()
		},
	})

	builder.Add(&di.Def{
		Name: "app",
		Build: func(ctn di.Container) (interface{}, error) {
			return createApplication(ctn), nil
		},
	})

	ctn, _ := builder.Build()
	defer ctn.DeleteWithSubContainers()

	var wg sync.WaitGroup
	discordReady := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		discord := ctn.Get("discord_service").(*DiscordService)
		discord.Start()
		close(discordReady)
	}()

	wailsApp := ctn.Get("app").(*application.App)
	createMainWindow(wailsApp)

	<-discordReady

	wailsApp.Events.On("discord_check", func(ev *application.WailsEvent) {
		discord := ctn.Get("discord_service").(*DiscordService)
		if discord.session.State.Ready.SessionID != "" {
			wailsApp.Events.Emit(&application.WailsEvent{Name: "discord_ready", Sender: "App", Data: map[string]interface{}{
				"name":     discord.session.State.Ready.User.Username,
				"icon_url": discord.session.State.Ready.User.AvatarURL("64x64"),
			}})
		}
	})
	err = wailsApp.Run()

	if err != nil {
		ctn.Get("logger").(*slog.Logger).Error("Failed to run application", "error", err)
		os.Exit(1)
	}

	wg.Wait()
}

func createApplication(ctn di.Container) *application.App {
	logger := ctn.Get("logger").(*slog.Logger)
	return application.New(application.Options{
		Name: "Jamkstrecka",
		Assets: application.AssetOptions{
			Handler:        application.AssetFileServerFS(Assets),
			DisableLogging: true,
		},
		Logger: logger.With("service", "APP"),
		Services: []application.Service{
			application.NewService(transactions.New(ctn)),
		},
		OnShutdown: func() {
			logger.Info("Shutting down application...")
		},
	})
}

func createMainWindow(app *application.App) {
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:             "Main Window",
		Width:            1024,
		Height:           768,
		Title:            "Jamkstrecka",
		URL:              "/",
		BackgroundColour: application.NewRGB(27, 38, 54),
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
	})
}

type ContainerAdapter struct {
	container di.Container
}

func (a *ContainerAdapter) Get(key string) interface{} {
	return a.container.Get(key)
}

type DiscordService struct {
	container di.Container
	logger    *slog.Logger
	session   *discordgo.Session
}

func NewDiscordService(ctn di.Container) *DiscordService {
	logger := ctn.Get("logger").(*slog.Logger)
	config := ctn.Get("config").(env.Config)
	logger.Info("Discord service starting...")

	session, err := discordgo.New("Bot " + config.DiscordToken)
	session.Identify.Intents |= discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	if err != nil {
		logger.Error("Failed to create discord session", "error", err)
		return nil
	}

	return &DiscordService{
		logger:    logger.With("service", "DISCORD"),
		session:   session,
		container: ctn,
	}
}

func (d *DiscordService) Start() {
	k, err := ken.New(d.session, ken.Options{
		DependencyProvider: &ContainerAdapter{container: d.container},
	})
	if err != nil {
		d.logger.Error("Failed to create ken", "error", err)
		return
	}

	err = k.RegisterCommands(
		new(commands.HelpCommand),
		new(commands.UserCommand),
		new(commands.StreckaCommand),
		new(commands.ProductCommand),
		new(commands.BalanceCommand),
		new(commands.PrintCommand),
	)

	if err != nil {
		d.logger.Error("Failed to register commands", "error", err)
		return
	}

	err = d.session.Open()
	if err != nil {
		d.logger.Error("Failed to open discord session", "error", err)
		return
	}

	d.logger.Info("Discord service started", "bot_name", d.session.State.Ready.User.Username)
	wails := d.container.Get("app").(*application.App)
	wails.Events.Emit(&application.WailsEvent{Name: "discord_ready", Sender: "Discord", Data: map[string]interface{}{
		"name":     d.session.State.Ready.User.Username,
		"icon_url": d.session.State.Ready.User.AvatarURL("64x64"),
	}})
}

func NewLogger() *slog.Logger {
	w := os.Stdout

	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
		AddSource:  true,
	}

	handler := utils.NewSourceHandler(w, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
