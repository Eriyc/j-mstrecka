package commands

import (
	"fmt"
	"gostrecka/internal/utils/static"
	"gostrecka/services/database"
	"gostrecka/services/discord"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type ProductCommand struct{}

var (
	_ ken.SlashCommand        = (*ProductCommand)(nil)
	_ ken.AutocompleteCommand = (*ProductCommand)(nil)
)

func (c *ProductCommand) Name() string {
	return "product"
}

func (c *ProductCommand) Description() string {
	return "Kommando för produkter"
}

func (c *ProductCommand) Version() string {
	return "1.0.0"
}

func (c *ProductCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *ProductCommand) Options() []*discordgo.ApplicationCommandOption {
	var integerOptionMinValue float64 = 1.0

	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Skapar en produkt",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Name of product",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "purchase_price",
					Description: "Vad betalade du för den när den köptes in?",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "internal_price",
					Description: "Internpris?",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "external_price",
					Description: "Externpris?",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "stock",
			Description: "Lägg till lagersaldo",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "product",
					Description:  "Produkt att updatera lagersaldo för",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "Antal att lägga till",
					Required:    true,
					MinValue:    &integerOptionMinValue,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Användare att lägga till för",
					Required:    false,
				},

				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "purchase_price",
					Description: "Uppdatera inköpspris?",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "internal_price",
					Description: "Uppdatera internpris?",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "external_price",
					Description: "Uppdatera externpris?",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "info",
			Description: "Visa information om en produkt",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "product",
					Description:  "Produkt att visa information om",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}
}

func (c *ProductCommand) IsDmCapable() bool {
	return true
}

func (c *ProductCommand) Autocomplete(ctx *ken.AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	return discord.AutocompleteSubcommand(ctx)
}

func (c *ProductCommand) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "create", Run: c.create},
		ken.SubCommandHandler{Name: "stock", Run: c.stock},
		ken.SubCommandHandler{Name: "info", Run: c.info},
	)

	return
}

func (c *ProductCommand) info(ctx ken.SubCommandContext) (err error) {
	productArg := ctx.Options().GetByName("product")

	log.Printf("productArg: %v", productArg.StringValue())

	ProductID, err := strconv.ParseInt(productArg.StringValue(), 10, 64)
	if err != nil {
		log.Printf("error converting product Id to int64: %v", err)
		return ctx.RespondError("Intern fel", "Fel")
	}

	db := ctx.Get(static.DiDatabase).(database.Database)
	product, price, err := db.GetProductIdent(ProductID)
	if err != nil {
		fmt.Printf("error getting product: %v", err)
		return ctx.RespondError("Produkten hittades inte", "Fel")
	}

	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Produkt",
		Description: product.Name,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Inköpspris",
				Value:  fmt.Sprintf("%.2f", price.PurchasePrice),
				Inline: true,
			},
			{
				Name:   "Internpris",
				Value:  fmt.Sprintf("%.2f", price.InternalPrice),
				Inline: true,
			},
			{
				Name:   "Externpris",
				Value:  fmt.Sprintf("%.2f", price.ExternalPrice),
				Inline: true,
			},

			{
				Name:   "Lagersaldo",
				Value:  fmt.Sprintf("%dst", product.TotalStock),
				Inline: true,
			},
		},
	})

	return
}

func (c *ProductCommand) create(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}
	name := ctx.Options().GetByName("name").StringValue()
	purchasePrice := ctx.Options().GetByName("purchase_price").FloatValue()
	internalPrice := ctx.Options().GetByName("internal_price").FloatValue()
	externalPrice := ctx.Options().GetByName("external_price").FloatValue()

	db := ctx.Get(static.DiDatabase).(database.Database)

	err = db.CreateProduct(name, purchasePrice, internalPrice, externalPrice)

	if err != nil {
		return nil
	}

	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Produkt",
		Description: "Produkt skapad: " + name,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Dags att lägga till lagersaldo!",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Inköpspris",
				Value:  fmt.Sprintf("%.2f", purchasePrice),
				Inline: true,
			},
			{
				Name:   "Internpris",
				Value:  fmt.Sprintf("%.2f", internalPrice),
				Inline: true,
			},
			{
				Name:   "Externpris",
				Value:  fmt.Sprintf("%.2f", externalPrice),
				Inline: true,
			},
		},
	})

	return
}

func (c *ProductCommand) stock(ctx ken.SubCommandContext) (err error) {
	productArg := ctx.Options().GetByName("product")
	amount := ctx.Options().GetByName("amount")
	userArg, userSupplied := ctx.Options().GetByNameOptional("user")

	purchasePrice, ppExists := ctx.Options().GetByNameOptional("purchase_price")
	internalPrice, ipExists := ctx.Options().GetByNameOptional("internal_price")
	externalPrice, epExists := ctx.Options().GetByNameOptional("external_price")

	db := ctx.Get(static.DiDatabase).(database.Database)

	var discordUser *discordgo.User
	if userSupplied {
		discordUser = userArg.UserValue(ctx)
	} else {
		discordUser = ctx.User()
	}

	user, oldWallet, err := db.GetUser(discordUser.ID)
	if err != nil {
		return ctx.RespondError("Användaren är inte registrerad i systemet, registrera med /user create <person>", "Fel")
	}

	ProductID, err := strconv.ParseInt(productArg.StringValue(), 10, 64)
	if err != nil {
		log.Printf("error converting product Id to int64: %v", err)
		return ctx.RespondError("Intern fel", "Fel")
	}

	product, price, err := db.GetProductIdent(ProductID)
	if err != nil {
		fmt.Printf("error getting product: %v", err)
		return ctx.RespondError("Produkten hittades inte", "Fel")
	}

	err = db.AddStock(product.ID, user.ID, amount.IntValue())
	if err != nil {
		fmt.Printf("error getting product: %v", err)
		return ctx.RespondError("Kunde inte lägga till lagersaldo", "Fel")
	}

	if ppExists || ipExists || epExists {
		if ppExists {
			price.PurchasePrice = purchasePrice.FloatValue()
		}
		if ipExists {
			price.InternalPrice = internalPrice.FloatValue()
		}
		if epExists {
			price.ExternalPrice = externalPrice.FloatValue()
		}

		err = db.UpdatePrice(product.ID, price.PurchasePrice, price.InternalPrice, price.ExternalPrice)
		if err != nil {
			return ctx.RespondError("Kunde inte uppdatera pris", "Fel")
		}
	}

	user, wallet, err := db.GetUser(discordUser.ID)
	if err != nil {
		return ctx.RespondError("Kunde inte hämta användare", "Fel")
	}

	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Lagersaldo",
		Description: fmt.Sprintf("Lagersaldo uppdaterat för %s", product.Name),

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Inköpare",
				Value: fmt.Sprintf("%s fick %.02f kronor att handla för 😋", discordUser.Mention(), wallet.TotalCreditsEarned-oldWallet.TotalCreditsEarned),
			},
			{
				Name:   "Inköpspris",
				Value:  fmt.Sprintf("%.2f", price.PurchasePrice),
				Inline: true,
			},
			{
				Name:   "Internpris",
				Value:  fmt.Sprintf("%.2f", price.InternalPrice),
				Inline: true,
			},
			{
				Name:   "Externpris",
				Value:  fmt.Sprintf("%.2f", price.ExternalPrice),
				Inline: true,
			},
		},
	})

	return nil
}
