package commands

import (
	"fmt"
	"gostrecka/internal/utils/static"
	"gostrecka/services/database"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type BalanceCommand struct{}

var (
	_ ken.SlashCommand = (*BalanceCommand)(nil)
)

func (c *BalanceCommand) Name() string {
	return "balance"
}

func (c *BalanceCommand) Description() string {
	return "Visar ditt saldo"
}

func (c *BalanceCommand) Version() string {
	return "1.0.0"
}

func (c *BalanceCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "Användaren att visa saldo för",
			Required:    false,
		},
	}
}

func (c *BalanceCommand) Run(ctx ken.Context) (err error) {
	userArg, ok := ctx.Options().GetByNameOptional("user")

	var selectedUser *discordgo.User
	if !ok {
		selectedUser = ctx.User()
	} else {
		selectedUser = userArg.UserValue(ctx)
	}

	db := ctx.Get(static.DiDatabase).(database.Database)

	user, balance, err := db.GetUser(selectedUser.ID)

	if err != nil {
		ctx.FollowUpError("Användaren finns inte", "")
	}

	var total string
	if balance.TotalCreditsEarned > balance.TotalDebtIncurred {
		total = fmt.Sprintf("%.02fkr i kredit", balance.TotalCreditsEarned-balance.TotalDebtIncurred)
	} else {
		total = fmt.Sprintf("%.02fkr i skuld", balance.TotalDebtIncurred-balance.TotalCreditsEarned)
	}

	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Saldo",
		Description: fmt.Sprintf("Saldo för <@%s>", user.ID),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Nuvarande",
				Value: total,
			},
			{
				Name:   "Total skuld",
				Value:  fmt.Sprintf("%.02fkr", balance.TotalDebtIncurred),
				Inline: true,
			},
			{
				Name:   "Totalt saldo",
				Value:  fmt.Sprintf("%.02fkr", balance.TotalCreditsEarned),
				Inline: true,
			},
		},
	})

	return
}
