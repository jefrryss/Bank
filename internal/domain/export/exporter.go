package export

import "task2/internal/domain/entities"

type Exporter interface {
	ExportBankAccount(account *[]entities.BankAccount)
	ExportCategory(category *[]entities.Category)
	ExportOperation(operation *[]entities.Operation)
	GetPath() error
}
