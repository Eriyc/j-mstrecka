package inits

import (
	"gostrecka/internal/utils/env"
	"gostrecka/internal/utils/snowflakenodes"
	"gostrecka/internal/utils/static"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
)

func InitDiscordBotSession(container di.Container) (release func()) {
	release = func() {}

	cfg := container.Get(static.DiConfig).(env.Config)

	if cfg.DiscordToken == "" {
		log.Printf("Discord token is not set")
		return
	}

	err := snowflakenodes.Setup()
	if err != nil {
		log.Fatalf("Failed setting up snowflake nodes %v", err)
	}

	session := container.Get(static.DiDiscordSession).(*discordgo.Session)

	session.Token = "Bot " + cfg.DiscordToken
	session.Identify.Intents |= discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	err = session.Open()
	if err != nil {
		log.Fatalf("Error opening DISCORD session %v", err)
	}

	return
}
