package main

import (
	"context"
	"embed"
	"fmt"
	inits "gostrecka/internal/init"
	"gostrecka/internal/service/database"
	"gostrecka/internal/service/desktop"
	"gostrecka/internal/utils/env"
	"gostrecka/internal/utils/static"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/zekrotja/ken"
)

//go:embed all:frontend/dist
var Assets embed.FS

func main() {
	desktop.Embeds = Assets

	store, err := env.NewConfigStore()
	if err != nil {
		fmt.Printf("could not initialize the config store: %v\n", err)
		return
	}

	fmt.Println(store.ConfigPath)

	cfg, err := store.Config()
	if err != nil {
		fmt.Printf("could not retrieve the configuration: %v\n", err)
		return
	}

	builder, _ := di.NewEnhancedBuilder()

	builder.Add(&di.Def{Name: static.DiConfig, Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }})

	builder.Add(&di.Def{
		Name: static.DiDatabase,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitDatabase(ctn), nil
		},
		Close: func(obj interface{}) error {
			database := obj.(database.Database)
			log.Println("Shutting down database connection...")
			database.Close()
			return nil
		},
	})

	builder.Add(&di.Def{
		Name: static.DiDiscordSession,
		Build: func(ctn di.Container) (interface{}, error) {
			return discordgo.New("")
		},
		Close: func(obj interface{}) error {
			session := obj.(*discordgo.Session)
			log.Println("Shutting down bot session...")
			session.Close()
			return nil
		},
	})

	// Initialize command handler
	builder.Add(&di.Def{
		Name: static.DiCommandHandler,
		Build: func(ctn di.Container) (interface{}, error) {
			return inits.InitCommandHandler(ctn)
		},
		Close: func(obj interface{}) error {
			log.Println("Unegister commands ...")
			return obj.(*ken.Ken).Unregister()
		},
	})

	builder.Add(&di.Def{
		Name: static.DiDesktop,
		Build: func(ctn di.Container) (interface{}, error) {
			return desktop.New(ctn)
		},
	})

	ctn, _ := builder.Build()
	defer ctn.DeleteWithSubContainers()

	ctn.Get(static.DiCommandHandler)
	releaseShard := inits.InitDiscordBotSession(ctn)
	defer releaseShard()

	InitDesktop(ctn)
}

func InitDesktop(container di.Container) (app *desktop.App) {
	app = container.Get(static.DiDesktop).(*desktop.App)

	err := wails.Run(&options.App{
		Title:    "gostrecka",
		Width:    1024,
		Height:   768,
		Logger:   logger.NewDefaultLogger(),
		LogLevel: logger.INFO,
		AssetServer: &assetserver.Options{
			Assets: Assets,
		},
		OnShutdown: func(ctx context.Context) {
			log.Println("Shutting down...")
			container.Delete()
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnDomReady:       app.OnDomReady,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return
}
