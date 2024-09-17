package commands

import (
	"fmt"
	"gostrecka/internal/utils/static"
	"gostrecka/services/database"
	"gostrecka/services/discord"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/zekrotja/ken"
)

type StreckaCommand struct{}

var (
	_ ken.SlashCommand        = (*StreckaCommand)(nil)
	_ ken.DmCapable           = (*StreckaCommand)(nil)
	_ ken.AutocompleteCommand = (*StreckaCommand)(nil)
)

func (c *StreckaCommand) Name() string {
	return "strecka"
}

func (c *StreckaCommand) Description() string {
	return "Streckar en produkt åt dig (eller någon annan)"
}

func (c *StreckaCommand) Version() string {
	return "1.0.0"
}

func (c *StreckaCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *StreckaCommand) Options() []*discordgo.ApplicationCommandOption {
	var integerOptionMinValue float64 = 1.0

	return []*discordgo.ApplicationCommandOption{
		{
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "product",
			Description:  "Produkten att strecka",
			Required:     true,
			Autocomplete: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "Användaren att strecka åt",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "amount",
			Description: "Antal att strecka",
			Required:    false,
			MinValue:    &integerOptionMinValue,
		},
	}
}

func (c *StreckaCommand) Autocomplete(ctx *ken.AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	return discord.AutocompleteOption(ctx)
}

func (c *StreckaCommand) IsDmCapable() bool {
	return true
}

func (c *StreckaCommand) Run(ctx ken.Context) (err error) {
	productArg := ctx.Options().GetByName("product")
	user, userSupplied := ctx.Options().GetByNameOptional("user")
	amountArg, amountSupplied := ctx.Options().GetByNameOptional("amount")

	var amount int64
	var response = "Streckar "
	if amountSupplied {
		amount = amountArg.IntValue()
		response += fmt.Sprintf("%vst ", amountArg.IntValue())
	} else {
		response += "1st "
		amount = 1
	}

	ProductID, err := strconv.ParseInt(productArg.StringValue(), 10, 64)
	if err != nil {
		log.Printf("error converting product Id to int64: %v", err)
		return ctx.RespondError("Intern fel", "Fel")
	}

	product, price, err := ctx.Get("database").(database.Database).GetProductIdent(ProductID)
	if err != nil {
		fmt.Printf("error getting product: %v", err)
		return ctx.RespondError("Produkten hittades inte", "Fel")
	}

	response += fmt.Sprintf("%s (%.02f)", product.Name, price.InternalPrice*float64(amount))

	var discordUser *discordgo.User
	if userSupplied {
		discordUser = user.UserValue(ctx)
		response += fmt.Sprintf(" åt %s", discordUser.Mention())
	} else {
		discordUser = ctx.User()
		response += " åt dig"
	}

	db := ctx.Get("database").(database.Database)
	userStruct, _, err := db.GetUser(discordUser.ID)

	if err != nil {
		ctx.Respond(&discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Content: "Du är inte registrerad i systemet\nRegistrera med /user create"}})
		return
	}

	_ = db.Strecka(userStruct, ProductID, amount)

	ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})

	desktop := ctx.Get("app").(*application.App)
	desktop.Events.Emit(&application.WailsEvent{Name: "transaction_updated", Sender: static.DiDesktop})

	return
}
