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
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/zekrotja/ken"
)

//go:embed all:frontend/dist
var Assets embed.FS

type SharedResources struct {
	DB     *Database
	Config env.Config
	Logger *slog.Logger
}

type AppContext struct {
	WailsApp *application.App
	Discord  *DiscordService
	Shared   *SharedResources
}

func main() {
	runtime.LockOSThread()

	shared := initSharedResources()

	var wg sync.WaitGroup
	discordReady := make(chan struct{})

	appCtx := &AppContext{
		Shared: shared,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		discord := NewDiscordService(shared, appCtx)
		appCtx.Discord = discord
		discord.Start()
		close(discordReady)
	}()

	wailsApp := createApplication(shared, appCtx)
	appCtx.WailsApp = wailsApp
	createMainWindow(wailsApp)

	<-discordReady

	err := wailsApp.Run()
	if err != nil {
		shared.Logger.Error("Failed to run application", "error", err)
		os.Exit(1)
	}

	wg.Wait()
}

func initSharedResources() *SharedResources {
	logger := NewLogger()

	configStore, err := env.NewConfigStore(logger.With("service", "CONFIG"))
	if err != nil {
		logger.Error("Failed to create config store", "error", err, "service", "CONFIG")
		os.Exit(1)
	}

	cfg, err := configStore.Config()
	if err != nil {
		configStore.Logger.Error("Failed to read config", "error", err)
		os.Exit(1)
	}

	return &SharedResources{
		DB:     NewDatabase(logger),
		Config: cfg,
		Logger: logger,
	}
}

func createApplication(shared *SharedResources, appCtx *AppContext) *application.App {
	return application.New(application.Options{
		Name: "Jamkstrecka",
		Assets: application.AssetOptions{
			Handler:        application.AssetFileServerFS(Assets),
			DisableLogging: true,
		},
		Logger: shared.Logger.With("service", "APP"),
		OnShutdown: func() {
			shared.Logger.Info("Shutting down application...")
			appCtx.Discord.session.Close()
		},
		// You can add services here if needed
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
	shared  *SharedResources
	appCtx  *AppContext
	logger  *slog.Logger
	session *discordgo.Session
}

func NewDiscordService(shared *SharedResources, appCtx *AppContext) *DiscordService {
	shared.Logger.Info("Discord service starting...")

	session, err := discordgo.New("Bot " + shared.Config.DiscordToken)
	if err != nil {
		shared.Logger.Error("Failed to create discord session", "error", err)
		return nil
	}

	return &DiscordService{
		shared:  shared,
		appCtx:  appCtx,
		logger:  shared.Logger.With("service", "DISCORD"),
		session: session,
	}
}

func (d *DiscordService) Get(key string) interface{} {
	// help
	return interface{}
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
