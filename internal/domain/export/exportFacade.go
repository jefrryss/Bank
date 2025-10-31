package export

import (
	"fmt"
	"task2/internal/domain/entities"
)

type ExportFacade struct {
	CommandsExport []Visitable
}

func NewExportFacade() *ExportFacade {
	return &ExportFacade{
		CommandsExport: make([]Visitable, 0),
	}
}
func (exporter *ExportFacade) Init(categories *[]*entities.Category, operations *[]*entities.Operation, accounts *[]*entities.BankAccount) {

	exporter.CommandsExport = []Visitable{}

	categoryAdapter := &CategoryAdapter{Categorys: categories}
	operationAdapter := &OperationAdapter{Operations: operations}
	accountAdapter := &BankAccountAdapter{BankAccounts: accounts}

	exporter.CommandsExport = append(exporter.CommandsExport, categoryAdapter, operationAdapter, accountAdapter)

}

func (e *ExportFacade) StartExport(typeExport ExporterVisitor) error {
	err := typeExport.GetPath()
	if err != nil {
		return err
	}
	for _, exportType := range e.CommandsExport {
		err := exportType.Accept(typeExport)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
