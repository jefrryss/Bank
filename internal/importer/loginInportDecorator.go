package importer

import (
	"fmt"
	"task2/domain/entities"
	"task2/domain/errordata"
	"task2/domain/ports"
	"task2/internal/logstorage"
	"time"
)

type LoggingImportDecorator struct {
	Importer ports.DataImporter
	LogStore *logstorage.LogStorage
}

func (d *LoggingImportDecorator) GetData() (
	[]entities.BankAccount,
	[]entities.Category,
	[]entities.Operation,
	[]errordata.ErrorRecord,
	error,
) {
	start := time.Now()

	accounts, categories, operations, errorData, err := GetData(d.Importer)
	elapsed := time.Since(start)

	src := fmt.Sprintf("%T", d.Importer)
	if csvImp, ok := d.Importer.(*ImportCSV); ok {
		src = fmt.Sprintf("CSV (%s)", csvImp.filePath)
	} else if jsonImp, ok := d.Importer.(*ImportJSON); ok {
		src = fmt.Sprintf("JSON (%s)", jsonImp.filePath)
	}

	accN := len(accounts)
	catN := len(categories)
	opN := len(operations)
	errN := len(errorData)

	startStr := start.Local().Format("2006-01-02 15:04:05")
	endStr := start.Add(elapsed).Local().Format("2006-01-02 15:04:05")
	ms := elapsed.Milliseconds()
	humanDur := elapsed.Truncate(time.Millisecond).String()
	if elapsed < time.Second {
		humanDur = fmt.Sprintf("%dms", ms)
	}

	title := fmt.Sprintf(
		"[%s → %s | %s/%d ms] Импорт из %s — accounts:%d, categories:%d, operations:%d, \n errors:%d",
		startStr, endStr, humanDur, ms, src, accN, catN, opN, errN,
	)

	var status string
	if err != nil {
		status = fmt.Sprintf("Ошибка: %s", err.Error())
	} else if errN > 0 {
		status = fmt.Sprintf("Завершено с частичными ошибками (%d)", errN)
	}

	d.LogStore.Add(logstorage.LogRecord{
		CommandName: title,
		StartedAt:   start,
		Duration:    elapsed,
		ErrorText:   status,
	})

	return accounts, categories, operations, errorData, err
}
