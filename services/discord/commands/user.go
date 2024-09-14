package commands

import (
	"database/sql"
	"fmt"
	"gostrecka/internal/service/database"
	"gostrecka/internal/utils/static"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type UserCommand struct{}

var (
	_ ken.SlashCommand = (*UserCommand)(nil)
	_ ken.DmCapable    = (*UserCommand)(nil)
)

func (c *UserCommand) Name() string {
	return "user"
}

func (c *UserCommand) Description() string {
	return "Kommando för användare"
}

func (c *UserCommand) Version() string {
	return "1.0.0"
}

func (c *UserCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *UserCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Skapar ett konto åt dig",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionUser,
					Name:         "user",
					Description:  "User to create for",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
	}
}

func (c *UserCommand) IsDmCapable() bool {
	return true
}

func (c *UserCommand) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "create", Run: c.create},
	)

	return
}

func (c *UserCommand) create(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}
	user, exists := ctx.Options().GetByNameOptional("user")

	var account *discordgo.User

	if exists {
		account = user.UserValue(ctx)
	} else {
		account = ctx.User()

	}
	db := ctx.Get(static.DiDatabase).(database.Database)

	match, _, err := db.GetUser(account.ID)
	if err != nil || match.ID != "" {
		log.Printf("Error: %v", err)
		if err != sql.ErrNoRows {
			err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Description: "User already exists",
			}).Send().Error

			if err != nil {
				panic(err)
			}

			return
		}
	}

	err = db.CreateUser(account.ID, account.GlobalName)

	if err != nil {
		err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Failed to create user: %s", err.Error()),
			Color:       1,
		}).Send().Error
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("User %s created!", account.Mention()),
	}).Send().Error
	return
}
