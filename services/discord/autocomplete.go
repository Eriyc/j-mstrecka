package discord

import (
	"fmt"
	"gostrecka/internal/utils/static"
	"gostrecka/services/database"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

func autocompleteInner(ctx *ken.AutocompleteContext, input string) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	input = strings.ToLower(input)

	db := ctx.Get(static.DiDatabase).(database.Database)
	items, err := db.SearchProduct(input)

	if err != nil {
		return nil, err
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(items))
	for _, item := range items {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  item.Product.Name,
			Value: fmt.Sprintf("%d", item.Product.ID),
		})
	}

	return choices, nil
}

func AutocompleteSubcommand(ctx *ken.AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	inputArg, ok := ctx.SubCommand().GetInput("product")

	var input = ""
	if !ok {
		log.Println("Could not get 'product'")
		return nil, nil
	} else {
		input = inputArg
	}

	return autocompleteInner(ctx, input)
}

func AutocompleteOption(ctx *ken.AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	inputArg, ok := ctx.GetInput("product")

	var input = ""
	if !ok {
		log.Println("Could not get 'product'")
		return nil, nil
	} else {
		input = inputArg
	}

	return autocompleteInner(ctx, input)
}
