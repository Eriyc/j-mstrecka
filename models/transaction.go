package models

import "time"

type Transaction struct {
	ID       int64   `json:"id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	UserName string  `json:"user_name"`
}

type LatestTransaction struct {
	UserID                     string    `json:"user_id"`
	UserName                   string    `json:"user_name"`
	TransactionDate            time.Time `json:"transaction_date"`
	CumulativeTransactionCount int64     `json:"cumulative_transaction_count"`
}

type TransactionLeaderboard struct {
	UserID                string `json:"user_id"`
	UserName              string `json:"user_name"`
	CurrentRank           int64  `json:"current_rank"`
	RankChangeIndicator   string `json:"rank_change_indicator"`
	TotalTransactionCount int64  `json:"total_transaction_count"`
}

type TransactionNumber struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int64   `json:"quantity"`
	PricePaid   float64 `json:"price_paid"`
}
