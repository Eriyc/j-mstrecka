package commands

import (
	"fmt"
	"gostrecka/services/database"
	"gostrecka/utils"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PrintCommand struct{}

// Description implements ken.SlashCommand.
func (p *PrintCommand) Description() string {
	return "Returns pdfs' of scannable barcodes"
}

// Name implements ken.SlashCommand.
func (p *PrintCommand) Name() string {
	return "print"
}

// Options implements ken.SlashCommand.
func (p *PrintCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

// Run implements ken.SlashCommand.
func (p *PrintCommand) Run(ctx ken.Context) (err error) {
	messageId := ctx.GetEvent().ID

	db := ctx.Get("database").(database.Database)

	upcRows, _ := db.GetUserUpcs()
	productRows, _ := db.GetProductUpcs()

	var users []utils.BarcodeInfo
	for _, upc := range upcRows {
		users = append(users, utils.BarcodeInfo{Number: upc.Upc, Label: upc.ReferableName})
	}

	var products []utils.BarcodeInfo
	for _, upc := range productRows {
		products = append(products, utils.BarcodeInfo{Number: upc.Upc, Label: upc.ReferableName})
	}

	err = utils.GenerateBarcodePDF(users, fmt.Sprintf("%s-users.pdf", messageId))
	if err != nil {
		return ctx.RespondError(err.Error(), "Error generating user barcode PDF")
	}
	err = utils.GenerateBarcodePDF(products, fmt.Sprintf("%s-products.pdf", messageId))
	if err != nil {
		return ctx.RespondError(err.Error(), "Error generating product barcode PDF")
	}

	userPDF, err := os.Open(fmt.Sprintf("%s-users.pdf", messageId))
	if err != nil {
		return err
	}
	defer userPDF.Close()

	productPDF, err := os.Open(fmt.Sprintf("%s-products.pdf", messageId))
	if err != nil {
		return err
	}
	defer productPDF.Close()

	ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Files: []*discordgo.File{
				{
					Name:        fmt.Sprintf("%s-users.pdf", messageId),
					ContentType: "application/pdf",
					Reader:      userPDF,
				},
				{
					Name:        fmt.Sprintf("%s-products.pdf", messageId),
					ContentType: "application/pdf",
					Reader:      productPDF,
				},
			},
		},
	})

	os.Remove(fmt.Sprintf("%s-users.pdf", messageId))
	os.Remove(fmt.Sprintf("%s-products.pdf", messageId))

	return
}

// Version implements ken.SlashCommand.
func (p *PrintCommand) Version() string {
	return "1.0.0"
}

var (
	_ ken.SlashCommand = (*PrintCommand)(nil)
)
