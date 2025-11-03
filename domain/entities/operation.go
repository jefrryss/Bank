package entities

import (
	"time"
)

type Operation struct {
	ID            string
	TypeOperation string
	Account       *BankAccount
	Amount        float64
	Date          time.Time
	Description   string
	CategoryID    *Category
}
