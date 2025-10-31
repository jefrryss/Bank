package analitic

import (
	"time"
)

type LogRecord struct {
	CommandName string
	StartedAt   time.Time
	Duration    time.Duration
	ErrorText   string
}
type LogStorage struct {
	records []LogRecord
}

func NewLogStorage() *LogStorage {
	return &LogStorage{records: []LogRecord{}}
}

func (ls *LogStorage) Add(record LogRecord) {
	ls.records = append(ls.records, record)
}

func (ls *LogStorage) GetAll() []LogRecord {
	return ls.records
}
