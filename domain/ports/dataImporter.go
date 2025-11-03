package ports

import (
	"github.com/jefrryss/Bank/domain/entities"
	"github.com/jefrryss/Bank/domain/errordata"
)

type DataImporter interface {
	ParseData() error
	SetFilePath(path string) error
	ParseBankAccounts() ([]entities.BankAccount, error)
	ParseCategories() ([]entities.Category, error)
	ParseOperations() ([]entities.Operation, error)
	GetErrorData() []errordata.ErrorRecord
}
