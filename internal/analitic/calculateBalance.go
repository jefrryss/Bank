package analitic

import (
	"fmt"

	"time"

	"github.com/jefrryss/Bank/domain/entities"
)

type CalculateBalanceCommand struct {
	Operations []entities.Operation
	From       time.Time
	To         time.Time
}

func NewCalculateBalanceCommand(ops []*entities.Operation) *CalculateBalanceCommand {
	operations := make([]entities.Operation, len(ops))
	for i, op := range ops {
		operations[i] = *op
	}

	return &CalculateBalanceCommand{
		Operations: operations,
	}
}

func (c *CalculateBalanceCommand) Execute() (string, error) {
	var income, expense float64

	for _, op := range c.Operations {
		if !op.Date.Before(c.From) && !op.Date.After(c.To) {
			if op.TypeOperation == "доход" {
				income += op.Amount
			} else {
				expense += op.Amount
			}
		}
	}

	result := fmt.Sprintf("Доход: %.2f | Расход: %.2f | Разница: %.2f",
		income, expense, income-expense)

	return result, nil
}
