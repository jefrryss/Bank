package export

import "task2/internal/domain/entities"

type ExporterVisitor interface {
	ExportBankAccount(account *[]*entities.BankAccount) error
	ExportCategory(category *[]*entities.Category) error
	ExportOperation(operation *[]*entities.Operation) error
	GetPath() error
}
