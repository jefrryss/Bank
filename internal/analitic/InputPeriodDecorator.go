package analitic

import (
	"fmt"
	"time"
)

// Декоратор над командой для добавления периода
type InputPeriodDecorator struct {
	Cmd  *CalculateBalanceCommand
	From time.Time
	To   time.Time
}

func (d *InputPeriodDecorator) Execute() (string, error) {
	d.Cmd.From = d.From
	d.Cmd.To = d.To

	result, err := d.Cmd.Execute()

	if err != nil {
		return "", err
	}

	periodInfo := fmt.Sprintf("Период: %s — %s\n\n",
		d.From.Format("2006-01-02"),
		d.To.Format("2006-01-02"))

	return periodInfo + result, nil
}
