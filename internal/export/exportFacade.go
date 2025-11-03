package export

import (
	"fmt"
	"task2/domain/entities"
	"task2/internal/logstorage"
)

// Фасад для работы с экспортом
type ExportFacade struct {
	CommandsExport []Visitable
	LogStorage     *logstorage.LogStorage
}

func NewExportFacade(logStorage *logstorage.LogStorage) *ExportFacade {
	return &ExportFacade{
		CommandsExport: make([]Visitable, 0),
		LogStorage:     logStorage,
	}
}
func (exporter *ExportFacade) Init(
	categories []*entities.Category,
	operations []*entities.Operation,
	accounts []*entities.BankAccount,
	logstorage *logstorage.LogStorage,
) {

	exporter.CommandsExport = []Visitable{}

	exporter.LogStorage = logstorage

	//Созадем декораторы для каждого адаптера, каждый адаптер будет посещать Visitor
	categoryDecorator := &LoginExportDecorator{exportAction: &CategoryAdapter{categories}, LogStore: exporter.LogStorage}
	operationDecorator := &LoginExportDecorator{exportAction: &OperationAdapter{operations}, LogStore: exporter.LogStorage}
	accountDecorator := &LoginExportDecorator{exportAction: &BankAccountAdapter{accounts}, LogStore: exporter.LogStorage}

	exporter.CommandsExport = append(exporter.CommandsExport, categoryDecorator, operationDecorator, accountDecorator)

}

func (e *ExportFacade) StartExport(typeExport ExporterVisitor) error {

	for _, exportType := range e.CommandsExport {
		err := exportType.Accept(typeExport)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
