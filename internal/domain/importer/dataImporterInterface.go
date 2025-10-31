package importer

import (
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
)

type DataImporter interface {
	ParseData() error
	GetPath() error
	ParseBankAccounts() ([]entities.BankAccount, error)
	ParseCategories() ([]entities.Category, error)
	ParseOperations() ([]entities.Operation, error)
	GetErrorData() []errordata.ErrorRecord
}

// Шаблоный метод (функция)
func GetValidateData(importer DataImporter) ([]entities.BankAccount, []entities.Category, []entities.Operation, []errordata.ErrorRecord, error) {

	err := importer.GetPath()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if err := importer.ParseData(); err != nil {
		return nil, nil, nil, nil, err
	}

	accounts, err := importer.ParseBankAccounts()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	categories, err := importer.ParseCategories()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	operations, err := importer.ParseOperations()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	errorData := importer.GetErrorData()
	return accounts, categories, operations, errorData, nil
}
