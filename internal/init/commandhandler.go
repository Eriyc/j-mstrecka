package inits

import (
	"gostrecka/internal/commands"
	"gostrecka/internal/utils/static"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/store"
)

// ContainerAdapter adapts di.Container to ken.ObjectProvider
type ContainerAdapter struct {
	container di.Container
}

// Get adapts the Get method to match ken.ObjectProvider's signature
func (a *ContainerAdapter) Get(key string) interface{} {
	return a.container.Get(key)
}

func InitCommandHandler(container di.Container) (k *ken.Ken, err error) {
	session := container.Get(static.DiDiscordSession).(*discordgo.Session)
	adapter := &ContainerAdapter{container: container}

	k, err = ken.New(session, ken.Options{
		CommandStore:       store.NewDefault(),
		DependencyProvider: adapter,
	})

	if err != nil {
		return
	}

	err = k.RegisterCommands(
		new(commands.HelpCommand),
		new(commands.UserCommand),
		new(commands.StreckaCommand),
		new(commands.ProductCommand),
		new(commands.BalanceCommand),
	)

	if err != nil {
		return
	}

	return
}
