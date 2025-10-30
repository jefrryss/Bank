package entities

import (
	"errors"
	"strings"
	"time"
)

type Operation struct {
	ID            string
	typeOperation string
	Account       *BankAccount
	Amount        float64
	Date          time.Time
	Description   string
	CategoryID    *Category
}

func (o *Operation) SetTypeOperation(typeInput string) error {
	typeInput = strings.ToLower(strings.TrimSpace(typeInput))

	switch typeInput {
	case "доход":
		o.typeOperation = "доход"
		return nil
	case "расход":
		o.typeOperation = "расход"
		return nil
	default:
		return errors.New("категория может быть только 'Доход' или 'Расход'")
	}
}

func (o *Operation) TypeOperation() string { return o.typeOperation }
