package desktop

import (
	"context"
	"fmt"
	"gostrecka/internal/utils/static"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/emersion/go-autostart"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type DiscordCheck struct {
	IconUrl string `json:"icon_url"`
	Name    string `json:"name"`
}

func (a *App) Startup(ctx context.Context) {
	a.Ctx = ctx

	// register autostart
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	auto := &autostart.App{
		Name: "jamkstrecka",
		Exec: []string{ex},
	}

	if !auto.IsEnabled() {
		err := auto.Enable()
		if err != nil {
			fmt.Printf("Could not enable autostart: %v", err)
			return
		}
	}

}

func (a *App) OnDomReady(ctx context.Context) {
	runtime.EventsOn(ctx, "discord_check", func(data ...interface{}) {
		session, _ := a.Container.Get(static.DiDiscordSession).(*discordgo.Session)
		if session.DataReady {
			var application = DiscordCheck{
				IconUrl: session.State.Ready.User.AvatarURL("64"),
				Name:    session.State.Ready.User.Username,
			}

			runtime.EventsEmit(ctx, "discord_ready", application)
		} else {
			runtime.EventsEmit(ctx, "discord_ready", nil)
		}
	})
}
