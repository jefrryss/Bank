package analitic

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Декоратор над командой для добавления вывода времени
type InputPeriodDecorator struct {
	Cmd *CalculateBalanceCommand
}

func (d *InputPeriodDecorator) Execute() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите дату ОТ (формат: 2006-01-02): ")
	fromStr, _ := reader.ReadString('\n')
	fromStr = strings.TrimSpace(fromStr)

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %v", err)
	}

	fmt.Print("Введите дату ДО (формат: 2006-01-02): ")
	toStr, _ := reader.ReadString('\n')
	toStr = strings.TrimSpace(toStr)

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %v", err)
	}

	d.Cmd.From = from
	d.Cmd.To = to

	fmt.Printf("Период установлен: %s — %s\n", from.Format("2006-01-02"), to.Format("2006-01-02"))

	return d.Cmd.Execute()
}
