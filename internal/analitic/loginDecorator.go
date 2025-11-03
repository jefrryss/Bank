package analitic

import (
	"fmt"
	"task2/internal/logstorage"
	"time"
)

// декоратор над командой для логирования
type LoggingDecoratorForAnalitic struct {
	Cmd      Command
	LogStore *logstorage.LogStorage
}

func (d *LoggingDecoratorForAnalitic) Execute() (string, error) {
	start := time.Now()

	out, err := d.Cmd.Execute()
	elapsed := time.Since(start)

	startStr := start.Local().Format("2006-01-02 15:04:05")
	endStr := start.Add(elapsed).Local().Format("2006-01-02 15:04:05")
	ms := elapsed.Milliseconds()
	humanDur := elapsed.Truncate(time.Millisecond).String()
	if elapsed < time.Second {
		humanDur = fmt.Sprintf("%dms", ms)
	}

	cmdType := fmt.Sprintf("%T", d.Cmd)

	title := fmt.Sprintf(
		"[%s → %s | %s/%d ms] Аналитика: %s",
		startStr, endStr, humanDur, ms, cmdType,
	)

	rec := logstorage.LogRecord{
		CommandName: title,
		StartedAt:   start,
		Duration:    elapsed,
	}
	if err != nil {
		rec.ErrorText = err.Error()
	}

	d.LogStore.Add(rec)
	return out, err
}
