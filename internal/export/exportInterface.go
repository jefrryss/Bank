package export

import "task2/domain/entities"

//Visitor

type ExporterVisitor interface {
	SetFilePath(path string) error
	ExportBankAccount(account []*entities.BankAccount) error
	ExportCategory(category []*entities.Category) error
	ExportOperation(operation []*entities.Operation) error
}
