package analitic

import (
	"fmt"
	"time"
)

type LoggingDecorator struct {
	Cmd      Command
	LogStore *LogStorage
}

func (d *LoggingDecorator) Execute() error {
	start := time.Now()

	err := d.Cmd.Execute()

	record := LogRecord{
		CommandName: fmt.Sprintf("%T", d.Cmd),
		StartedAt:   start,
		Duration:    time.Since(start),
	}

	if err != nil {
		record.ErrorText = err.Error()
	}

	d.LogStore.Add(record)

	return err
}
