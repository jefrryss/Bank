package export

import (
	"task2/internal/domain/entities"
)

type Visitable interface {
	Accept(exporter ExporterVisitor) error
}

type CategoryAdapter struct {
	Categorys *[]*entities.Category
}

func (cat *CategoryAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportCategory(cat.Categorys)
	return err
}

type BankAccountAdapter struct {
	BankAccounts *[]*entities.BankAccount
}

func (bank *BankAccountAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportBankAccount(bank.BankAccounts)
	return err
}

type OperationAdapter struct {
	Operations *[]*entities.Operation
}

func (oper *OperationAdapter) Accept(exporter ExporterVisitor) error {
	err := exporter.ExportOperation(oper.Operations)
	return err
}
