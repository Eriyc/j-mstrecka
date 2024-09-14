package desktop

import (
	"embed"
	"gostrecka/internal/models"
	"gostrecka/internal/service/database"
	"gostrecka/internal/utils/static"
	"log"
	"strconv"

	"github.com/sarulabs/di/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Desktop struct {
	Container di.Container
	App       application.App
	Window    application.Window
}

var Embeds embed.FS

func New(container di.Container) (desktop *Desktop, err error) {
	desktop = new(Desktop)
	desktop.Container = container

	desktop.App = *application.New(application.Options{
		Name:   "jamkstrecka",
		Assets: application.AlphaAssets,
	})

	return
}

func (a *Desktop) GetLatestTransactions() []models.LatestTransaction {
	db := a.Container.Get(static.DiDatabase).(database.Database)
	transactions, err := db.GetLatestTransactions()

	if err != nil {
		return []models.LatestTransaction{}
	}

	return transactions
}

func (a *Desktop) GetLeaderboard() []models.TransactionLeaderboard {
	db := a.Container.Get(static.DiDatabase).(database.Database)
	leaderboard, err := db.GetTransactionLeaderboard()

	if err != nil {
		return []models.TransactionLeaderboard{}
	}

	return leaderboard
}

func (a *Desktop) ScanUpc(upc string) interface{} {

	db := a.Container.Get(static.DiDatabase).(database.Database)
	log.Printf("scanning upc: %v", upc)
	result, err := db.GetUpcType(upc)
	if err != nil {
		log.Printf("error getting upc type: %v", err)
		return nil
	}

	if result.Type == "product" {
		id, _ := strconv.ParseInt(result.ReferableId, 10, 64)
		product, price, err := db.GetProductIdent(id)
		if err != nil {
			log.Printf("error getting product: %v", err)

			return nil
		}

		return map[string]interface{}{
			"type":    "product",
			"product": product,
			"price":   price,
		}
	}

	if result.Type == "user" {
		user, balance, err := db.GetUser(result.ReferableId)
		if err != nil {
			log.Printf("error getting user: %v", err)
			return nil
		}

		return map[string]interface{}{
			"type":    "user",
			"user":    user,
			"balance": balance,
		}
	}

	return nil
}

func (a *Desktop) Strecka(ProductID int64, UserID string, amount int64) (result interface{}) {
	db := a.Container.Get(static.DiDatabase).(database.Database)
	err := db.Strecka(models.User{ID: UserID}, ProductID, amount)

	if err != nil {
		log.Printf("error strecka: %v", err)
		return map[string]interface{}{
			"error":   err.Error(),
			"user":    nil,
			"product": nil,
			"balance": nil,
		}
	}

	user, balance, err := db.GetUser(UserID)
	if err != nil {
		log.Printf("error getting user: %v", err)
		return map[string]interface{}{
			"error":   err.Error(),
			"user":    nil,
			"product": nil,
			"balance": nil,
		}
	}

	product, _, err := db.GetProductIdent(ProductID)
	if err != nil {
		log.Printf("error getting product: %v", err)
		return map[string]interface{}{
			"error":   err.Error(),
			"user":    nil,
			"product": nil,
			"balance": nil,
		}
	}

	result = map[string]interface{}{
		"error":   nil,
		"user":    user,
		"product": product,
		"balance": balance,
	}

	a.App.Events.Emit(&application.WailsEvent{Name: "transaction_updated", Sender: static.DiDesktop})

	return
}
