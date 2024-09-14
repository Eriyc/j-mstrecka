package inits

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/fx"
)

func NewDesktop(lc fx.Lifecycle) *application.App {
	app := application.New(application.Options{Name: "Jamkstrecka"})

	return app
}

func RunDesktop(lc fx.Lifecycle, app *application.App) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {

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

				fmt.Println("Starting Wails aplication")
				app.Run()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
