package ports

import "task2/domain/entities"

//Visitor

type ExporterVisitor interface {
	SetFilePath(path string) error
	ExportBankAccount(account []*entities.BankAccount) error
	ExportCategory(category []*entities.Category) error
	ExportOperation(operation []*entities.Operation) error
}

// Элементы которые посещает Visitor(то есть реализации интерфейса ExportVisitor: ExportCSV, exportJSON)
type Visitable interface {
	Accept(exporter ExporterVisitor) error
}
