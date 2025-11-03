package analitic

import (
	"github.com/jefrryss/Bank/internal/logstorage"

	"github.com/jefrryss/Bank/domain/entities"

	"time"
)

type AnalyticsFacade struct {
	operations []*entities.Operation
	categories []*entities.Category
	invoker    *Invoker
	logStore   *logstorage.LogStorage
}

func NewAnalyticsFacade(ops []*entities.Operation, cat []*entities.Category, logstorage *logstorage.LogStorage) *AnalyticsFacade {
	return &AnalyticsFacade{
		operations: ops,
		categories: cat,
		invoker:    &Invoker{},
		logStore:   logstorage,
	}
}

func (f *AnalyticsFacade) CalcBalance(from, to time.Time) string {
	cmd := NewCalculateBalanceCommand(f.operations)

	cmdWithPeriod := &InputPeriodDecorator{
		Cmd:  cmd,
		From: from,
		To:   to,
	}

	loggedCmd := &LoggingDecoratorForAnalitic{
		Cmd:      cmdWithPeriod,
		LogStore: f.logStore,
	}

	f.invoker.Clear()
	f.invoker.AddCommand(loggedCmd)
	result, _ := f.invoker.Run()

	return result
}

func (f *AnalyticsFacade) GroupByCategory() string {
	cmd := NewGroupByCategoryCommand(f.operations, f.categories)

	loggedCmd := &LoggingDecoratorForAnalitic{
		Cmd:      cmd,
		LogStore: f.logStore,
	}

	f.invoker.Clear()
	f.invoker.AddCommand(loggedCmd)
	result, _ := f.invoker.Run()

	return result
}
