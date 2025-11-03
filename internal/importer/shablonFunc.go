package importer

import (
	"task2/domain/entities"
	"task2/domain/errordata"
	"task2/domain/ports"
)

// // Шаблоный метод (функция)
func GetData(importer ports.DataImporter) ([]entities.BankAccount, []entities.Category, []entities.Operation, []errordata.ErrorRecord, error) {

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
