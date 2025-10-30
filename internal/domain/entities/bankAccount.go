package entities

import (
	"errors"
)

type BankAccount struct {
	ID      string
	Name    string
	balance float64
}

func (b *BankAccount) SetBalance(balance float64) error {
	if balance < 0 {
		return errors.New("баланс не может быть отрицательным")
	}
	b.balance = balance
	return nil
}

func (b *BankAccount) Balance() float64 { return b.balance }
