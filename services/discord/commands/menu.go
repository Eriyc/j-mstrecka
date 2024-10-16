package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type MenuCommand struct{}

var (
	_ ken.SlashCommand = (*MenuCommand)(nil)
	_ ken.DmCapable    = (*MenuCommand)(nil)
)

func (c *MenuCommand) Name() string {
	return "menu"
}

func (c *MenuCommand) Description() string {
	return "shows whats avalible"
}

func (c *MenuCommand) Version() string {
	return "1.0.0"
}

func (c *MenuCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *MenuCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *MenuCommand) IsDmCapable() bool {
	return true
}

func (c *MenuCommand) Run(ctx ken.Context) (err error) {
	var embed = &discordgo.MessageEmbed{
		Title:       "JämK Meny",
		Description: "Detta kan du köpa",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "1 Miljon öl",
				Value: "10 miljon öl.",
			},
			{
				Name:  "500 cigaretter",
				Value: "Ät dem.",
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
