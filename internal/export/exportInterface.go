package export

import "github.com/jefrryss/Bank/domain/entities"

//Visitor

type ExporterVisitor interface {
	SetFilePath(path string) error
	ExportBankAccount(account []*entities.BankAccount) error
	ExportCategory(category []*entities.Category) error
	ExportOperation(operation []*entities.Operation) error
}
