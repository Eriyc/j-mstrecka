package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type HelpCommand struct{}

var (
	_ ken.SlashCommand = (*HelpCommand)(nil)
	_ ken.DmCapable    = (*HelpCommand)(nil)
)

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Basic Test Command"
}

func (c *HelpCommand) Version() string {
	return "1.0.0"
}

func (c *HelpCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *HelpCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *HelpCommand) IsDmCapable() bool {
	return true
}

func (c *HelpCommand) Run(ctx ken.Context) (err error) {

	var embed = &discordgo.MessageEmbed{
		Title:       "Hjälp",
		Description: "Jag är dum och kommer inte ihåg vad jag ska göra <:(",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "/user create [user]",
				Value:  "Registrerar ditt (eller någon annans) konto i systemet",
				Inline: true,
			},
			{
				Name:   "/balance [user]",
				Value:  "Visar hur mycket du (eller någon annan) har köpt för",
				Inline: false,
			},
			{
				Name:   "/strecka <product> [amount] [user]",
				Value:  "Streckar en produkt åt dig (eller någon annan)",
				Inline: false,
			},
		},
	}

	ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	return
}
