package main

import (
	"embed"
	"gostrecka/internal/utils/logger"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/lmittmann/tint"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var Assets embed.FS

type SharedResources struct {
	DB     *Database
	Config *Config
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
	return &SharedResources{
		DB:     NewDatabase(logger),
		Config: NewConfig(logger),
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
	shared *SharedResources
	appCtx *AppContext
	logger *slog.Logger
}

func NewDiscordService(shared *SharedResources, appCtx *AppContext) *DiscordService {
	return &DiscordService{
		shared: shared,
		appCtx: appCtx,
		logger: shared.Logger.With("service", "discord"),
	}
}

func (d *DiscordService) Start() {
	d.logger.Info("Discord service started")
	// Initialize and start your Discord bot here
}

type Database struct {
	logger *slog.Logger
}

func NewDatabase(logger *slog.Logger) *Database {
	return &Database{logger: logger.With("service", "DATABASE")}
}

type Config struct {
	logger *slog.Logger
}

func NewConfig(logger *slog.Logger) *Config {
	return &Config{logger: logger.With("service", "CONFIG")}
}

func NewLogger() *slog.Logger {
	w := os.Stdout

	opts := &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
		AddSource:  false,
	}

	handler := logger.NewSourceHandler(w, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
