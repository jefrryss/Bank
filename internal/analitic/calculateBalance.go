package analitic

import (
	"fmt"
	"task2/internal/domain/entities"
	"time"
)

// Команда для расчёта баланса за период
type CalculateBalanceCommand struct {
	Operations []entities.Operation
	From       time.Time
	To         time.Time
}

// Конструктор команды, принимает данные из BankManager
func NewCalculateBalanceCommand(ops []*entities.Operation) *CalculateBalanceCommand {
	operations := make([]entities.Operation, len(ops))
	for i, op := range ops {
		operations[i] = *op // разыменовываем указатель
	}

	return &CalculateBalanceCommand{
		Operations: operations,
	}
}

// Выполнение команды
func (c *CalculateBalanceCommand) Execute() error {
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

	fmt.Printf("Доход: %.2f | Расход: %.2f | Разница: %.2f\n",
		income, expense, income-expense)

	return nil
}
