package export

import (
	"github.com/jefrryss/Bank/domain/entities"
)

// Элементы которые посещает Visitor(реализации интерфейса ExportVisitor: ExportCSV, exportTJSON)
type Visitable interface {
	Accept(exporter ExporterVisitor) error
}

// Адаптер для категории
type CategoryAdapter struct {
	Categorys []*entities.Category
}

func (cat *CategoryAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportCategory(cat.Categorys)
	return err
}

// Адаптер для аккаунтов
type BankAccountAdapter struct {
	BankAccounts []*entities.BankAccount
}

func (bank *BankAccountAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportBankAccount(bank.BankAccounts)
	return err
}

// Адаптер для операций
type OperationAdapter struct {
	Operations []*entities.Operation
}

func (oper *OperationAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportOperation(oper.Operations)
	return err
}
