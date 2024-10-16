package database

import "gostrecka/models"

type Database interface {
	Connect() error
	Close()
	Status() error

	/* Users */
	GetUser(id string) (user models.User, balance models.Balance, err error)
	CreateUser(id string, name string) error

	/* Products */
	GetProductIdent(id int64) (product models.Product, price models.ProductPrice, err error)
	SearchProduct(name string) (products []models.ProductWithPrice, err error)
	CreateProduct(name string, purchasePrice float64, internalPrice float64, externalPrice float64) error

	UpdatePrice(productId int64, purchasePrice float64, internalPrice float64, externalPrice float64) error

	/* Stock */
	AddStock(productId int64, userId string, amount int64) error

	/* UPCs */
	GetUpcType(upc string) (lookup models.UpcLookup, err error)
	GetUserUpcs() (upcs []models.Upc, err error)
	GetProductUpcs() (upcs []models.Upc, err error)

	/* Transactions */
	Strecka(user models.User, productId int64, amount int64) error
	GetLatestTransactions() (transactions []models.LatestTransaction, err error)
	GetTransactionLeaderboard() (leaderboard []models.TransactionLeaderboard, err error)
	GetTransactionNumbers(UserId string) (items []models.TransactionNumber, err error)
}
