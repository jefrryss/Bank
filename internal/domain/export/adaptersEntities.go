package export

import "task2/internal/domain/entities"

type Visitable interface {
	Accpet(exporter Exporter)
}

type CategoryAdapter struct {
	Categorys *[]entities.Category
}

func (cat *CategoryAdapter) Accept(exporter Exporter) {
	exporter.ExportCategory(cat.Categorys)
}

type BankAccountAdapter struct {
	BankAccounts *[]entities.Category
}

func (bank *BankAccountAdapter) Accept(exporter Exporter) {
	exporter.ExportCategory(bank.BankAccounts)
}

type OperationAdapter struct {
	Operations *[]entities.Category
}

func (oper *OperationAdapter) Accept(exporter Exporter) {
	exporter.ExportCategory(oper.Operations)
}
