package models

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Balance struct {
	TotalCreditsEarned float64 `json:"total_credits_earned"`
	TotalPaymentsMade  float64 `json:"total_payments_made"`
	TotalDebtIncurred  float64 `json:"total_debt_incurred"`
	RemainingCredits   float64 `json:"remaining_credits"`
	DebtIncurred       float64 `json:"debt_incurred"`
}
