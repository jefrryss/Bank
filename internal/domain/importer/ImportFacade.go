package importer

import (
	"fmt"
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
	"task2/internal/domain/factory"
)

type ImportFacade struct {
	accountsValid      []entities.BankAccount
	operationValid     []entities.Operation
	categoryValid      []entities.Category
	operationFactory   factory.OperationFactory
	categoryFactory    factory.CategoryFactory
	bankAccountFactory factory.BankAccount
	errorData          []errordata.ErrorRecord
}

func NewImportFacade() *ImportFacade {
	return &ImportFacade{
		accountsValid:      make([]entities.BankAccount, 0),
		operationValid:     make([]entities.Operation, 0),
		categoryValid:      make([]entities.Category, 0),
		operationFactory:   factory.OperationFactory{},
		categoryFactory:    factory.CategoryFactory{},
		bankAccountFactory: factory.BankAccount{},
		errorData:          make([]errordata.ErrorRecord, 0),
	}
}
func (i *ImportFacade) Init(importer DataImporter) error {
	accounts, categories, operations, errorsList, err := GetValidateData(importer)
	if err != nil {
		return err
	}
	i.errorData = errorsList

	for _, acc := range accounts {
		b, err := i.bankAccountFactory.CreateBankAccount(acc.ID, acc.Name, acc.Balance)
		if err == nil {
			i.accountsValid = append(i.accountsValid, b)
		} else {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", acc),
				Err:  err,
			})
		}
	}
	for _, cat := range categories {
		c, err := i.categoryFactory.CreateCategory(
			cat.ID,
			cat.Name,
			cat.TypeCategory,
		)
		if err == nil {
			i.categoryValid = append(i.categoryValid, c)
		} else {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", cat),
				Err:  err,
			})
		}
	}

	for _, op := range operations {
		o, err := i.operationFactory.CreateOperation(
			op.ID,
			op.TypeOperation,
			op.Account,
			op.CategoryID,
			op.Amount,
			op.Date,
			op.Description,
		)
		if err == nil {
			i.operationValid = append(i.operationValid, o)
		} else {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", op),
				Err:  err,
			})
		}
	}

	return nil
}

func (i *ImportFacade) GetOperation() []entities.Operation {
	return i.operationValid
}

func (i *ImportFacade) GetCategory() []entities.Category {
	return i.categoryValid
}

func (i *ImportFacade) GetAccounts() []entities.BankAccount {
	return i.accountsValid
}

func (i *ImportFacade) GetErrors() []errordata.ErrorRecord {
	return i.errorData
}
