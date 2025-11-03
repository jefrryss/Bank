package importer

import (
	"github.com/jefrryss/Bank/domain/entities"
	"github.com/jefrryss/Bank/domain/errordata"
	"github.com/jefrryss/Bank/domain/ports"
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
