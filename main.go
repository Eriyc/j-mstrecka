package main

import (
	"embed"
	"gostrecka/internal/utils/logger"
	"gostrecka/services/discord/commands"
	"gostrecka/services/env"
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

//go:embed all:frontend/dist
var Assets embed.FS

func main() {
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
			logger := ctn.Get("logger").(*slog.Logger)
			return NewDatabase(logger), nil
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
		Name: "wails_app",
		Build: func(ctn di.Container) (interface{}, error) {
			return createApplication(ctn), nil
		},
	})

	if err != nil {
		panic(err)
	}

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

	wailsApp := ctn.Get("wails_app").(*application.App)
	createMainWindow(wailsApp)

	<-discordReady

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

type DiscordService struct {
	logger  *slog.Logger
	session *discordgo.Session
}

func NewDiscordService(ctn di.Container) *DiscordService {
	logger := ctn.Get("logger").(*slog.Logger)
	config := ctn.Get("config").(env.Config)
	logger.Info("Discord service starting...")

	session, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		logger.Error("Failed to create discord session", "error", err)
		return nil
	}

	return &DiscordService{
		logger:  logger.With("service", "DISCORD"),
		session: session,
	}
}

func (d *DiscordService) Start() {
	k, err := ken.New(d.session, ken.Options{})
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
}

type Database struct {
	logger *slog.Logger
}

func NewDatabase(logger *slog.Logger) *Database {
	return &Database{logger: logger.With("service", "DATABASE")}
}

func NewLogger() *slog.Logger {
	w := os.Stdout

	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
		AddSource:  false,
	}

	handler := logger.NewSourceHandler(w, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
