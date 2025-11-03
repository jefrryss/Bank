package ports

import (
	"task2/domain/entities"
	"task2/domain/errordata"
)

type DataImporter interface {
	ParseData() error
	SetFilePath(path string) error
	ParseBankAccounts() ([]entities.BankAccount, error)
	ParseCategories() ([]entities.Category, error)
	ParseOperations() ([]entities.Operation, error)
	GetErrorData() []errordata.ErrorRecord
}
