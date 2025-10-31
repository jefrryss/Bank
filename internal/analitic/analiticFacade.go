package analitic

import (
	"fmt"
	"task2/internal/domain/entities"
)

// AnalyticsFacade работает с данными напрямую (слайсы)
type AnalyticsFacade struct {
	operations []*entities.Operation
	categories []*entities.Category
	invoker    *Invoker
	logStore   *LogStorage
}

// Конструктор, получает данные от BankManager

func NewAnalyticsFacade(ops []*entities.Operation, cats []*entities.Category) *AnalyticsFacade {
	return &AnalyticsFacade{
		operations: ops,
		categories: cats,
		invoker:    &Invoker{},
		logStore:   NewLogStorage(),
	}
}

// Показать логи выполнения команд
func (f *AnalyticsFacade) ShowLogs() {
	fmt.Println("\n===== ЛОГИ ВЫПОЛНЕННЫХ КОМАНД =====")
	for _, r := range f.logStore.GetAll() {
		status := "успешно"
		if r.ErrorText != "" {
			status = "ошибка"
		}

		fmt.Printf("[%s] Команда: %s | Статус: %s | Длительность: %v\n",
			r.StartedAt.Format("2006-01-02 15:04:05"),
			r.CommandName,
			status,
			r.Duration,
		)

		if r.ErrorText != "" {
			fmt.Printf("Причина ошибки: %s\n", r.ErrorText)
		}
	}
	fmt.Println("===================================")
}

// Расчет баланса за период
func (f *AnalyticsFacade) CalcBalance() {
	cmd := NewCalculateBalanceCommand(f.operations)

	cmdWithPeriod := &InputPeriodDecorator{Cmd: cmd}

	loggedCmd := &LoggingDecorator{
		Cmd:      cmdWithPeriod,
		LogStore: f.logStore,
	}

	f.invoker.Clear()
	f.invoker.AddCommand(loggedCmd)
	f.invoker.Run()
}

// Группировка по категориям
func (f *AnalyticsFacade) GroupByCategory() {
	cmd := NewGroupByCategoryCommand(f.operations, f.categories)

	loggedCmd := &LoggingDecorator{
		Cmd:      cmd,
		LogStore: f.logStore,
	}

	f.invoker.Clear()
	f.invoker.AddCommand(loggedCmd)
	f.invoker.Run()
}
