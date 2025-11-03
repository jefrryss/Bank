package export

import (
	"fmt"
	"task2/internal/logstorage"
	"time"
)

// декоратор для экспорта, чтобы логи записывать
type LoginExportDecorator struct {
	exportAction Visitable //элемент в который заходит visitor
	LogStore     *logstorage.LogStorage
}

func (l *LoginExportDecorator) Accept(exporter ExporterVisitor) error {
	start := time.Now()

	err := l.exportAction.Accept(exporter)

	record := logstorage.LogRecord{
		CommandName: fmt.Sprintf("Тип экпорта:%T, Категория: %T", exporter, l.exportAction),
		StartedAt:   start,
		Duration:    time.Since(start),
	}

	if err != nil {
		record.ErrorText = err.Error()
	}

	l.LogStore.Add(record)

	return err
}
