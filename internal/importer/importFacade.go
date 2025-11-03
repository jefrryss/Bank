package importer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jefrryss/Bank/domain/errordata"
	"github.com/jefrryss/Bank/domain/factory"
	"github.com/jefrryss/Bank/domain/ports"
	"github.com/jefrryss/Bank/internal/logstorage"

	"github.com/jefrryss/Bank/domain/entities"
)

type ImportFacade struct {
	accountsValid      []entities.BankAccount
	operationValid     []entities.Operation
	categoryValid      []entities.Category
	operationFactory   factory.OperationFactory
	categoryFactory    factory.CategoryFactory
	bankAccountFactory factory.BankAccount
	errorData          []errordata.ErrorRecord
}

func NewImportFacade() *ImportFacade {
	return &ImportFacade{
		accountsValid:      make([]entities.BankAccount, 0),
		operationValid:     make([]entities.Operation, 0),
		categoryValid:      make([]entities.Category, 0),
		operationFactory:   factory.OperationFactory{},
		categoryFactory:    factory.CategoryFactory{},
		bankAccountFactory: factory.BankAccount{},
		errorData:          make([]errordata.ErrorRecord, 0),
	}
}
func (i *ImportFacade) Init(importer ports.DataImporter, logStore *logstorage.LogStorage) error {

	i.accountsValid = i.accountsValid[:0]
	i.operationValid = i.operationValid[:0]
	i.categoryValid = i.categoryValid[:0]
	i.errorData = i.errorData[:0]

	decorated := &LoggingImportDecorator{
		Importer: importer,
		LogStore: logStore,
	}
	accounts, categories, operations, errorsList, err := decorated.GetData()
	if err != nil {
		return err
	}
	if len(errorsList) > 0 {
		i.errorData = append(i.errorData, errorsList...)
	}

	for _, acc := range accounts {
		b, err := i.bankAccountFactory.CreateBankAccount(
			strings.TrimSpace(acc.ID),
			strings.TrimSpace(acc.Name),
			acc.Balance,
		)
		if err != nil {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf(
					"bank_account{id=%q,name=%q,balance=%v}",
					acc.ID, acc.Name, acc.Balance,
				),
				Err: fmt.Errorf("некорректный банк-аккаунт: %w", err),
			})
			continue
		}
		i.accountsValid = append(i.accountsValid, b)
	}

	for _, cat := range categories {
		c, err := i.categoryFactory.CreateCategory(
			strings.TrimSpace(cat.ID),
			strings.TrimSpace(cat.Name),
			strings.TrimSpace(cat.TypeCategory),
		)
		if err != nil {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf(
					"category{id=%q,name=%q,type=%q}",
					cat.ID, cat.Name, cat.TypeCategory,
				),
				Err: fmt.Errorf("некорректная категория: %w", err),
			})
			continue
		}
		i.categoryValid = append(i.categoryValid, c)
	}

	accIndex := make(map[string]*entities.BankAccount, len(i.accountsValid))
	for idx := range i.accountsValid {
		a := &i.accountsValid[idx]
		if a.ID != "" {
			accIndex[a.ID] = a
		}
	}
	catIndex := make(map[string]*entities.Category, len(i.categoryValid))
	for idx := range i.categoryValid {
		c := &i.categoryValid[idx]
		if c.ID != "" {
			catIndex[c.ID] = c
		}
	}

	for _, op := range operations {

		opLine := fmt.Sprintf(
			"operation{id=%q,type=%q,account_id=%q,amount=%v,date=%s,category_id=%q,desc=%q}",
			op.ID,
			op.TypeOperation,
			func() string {
				if op.Account != nil {
					return op.Account.ID
				}
				return ""
			}(),
			op.Amount,
			op.Date.Format("2006-01-02"),
			func() string {
				if op.CategoryID != nil {
					return op.CategoryID.ID
				}
				return ""
			}(),
			op.Description,
		)

		opID := strings.TrimSpace(op.ID)
		if opID == "" {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  errors.New("операция без id — пропущена"),
			})
			continue
		}

		opType := strings.ToLower(strings.TrimSpace(op.TypeOperation))
		if opType != "доход" && opType != "расход" {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: недопустимый type=%q (ожидалось \"доход\" или \"расход\")", opID, op.TypeOperation),
			})
			continue
		}

		var accID string
		if op.Account != nil {
			accID = strings.TrimSpace(op.Account.ID)
		}
		if accID == "" {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: не указан account_id (поле пустое)", opID),
			})
			continue
		}
		accRef, ok := accIndex[accID]
		if !ok {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: ссылочный account_id=%q отсутствует среди валидных аккаунтов", opID, accID),
			})
			continue
		}

		var catRef *entities.Category
		if op.CategoryID != nil && strings.TrimSpace(op.CategoryID.ID) != "" {
			catID := strings.TrimSpace(op.CategoryID.ID)
			c, ok := catIndex[catID]
			if !ok {
				i.errorData = append(i.errorData, errordata.ErrorRecord{
					Line: opLine,
					Err:  fmt.Errorf("операция id=%q: ссылочная категория id=%q не найдена среди валидных", opID, catID),
				})
				continue
			}
			catRef = c
		}

		if op.Amount <= 0 {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: amount=%v должен быть > 0", opID, op.Amount),
			})
			continue
		}
		if op.Date.IsZero() {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: дата отсутствует или некорректна", opID),
			})
			continue
		}

		o, err := i.operationFactory.CreateOperation(
			opID,
			opType,
			accRef,
			catRef,
			op.Amount,
			op.Date,
			strings.TrimSpace(op.Description),
		)
		if err != nil {
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: opLine,
				Err:  fmt.Errorf("операция id=%q: %w", opID, err),
			})
			continue
		}

		i.operationValid = append(i.operationValid, o)
	}

	return nil
}

func (i *ImportFacade) GetOperation() []entities.Operation {
	return i.operationValid
}

func (i *ImportFacade) GetCategory() []entities.Category {
	return i.categoryValid
}

func (i *ImportFacade) GetAccounts() []entities.BankAccount {
	return i.accountsValid
}

func (i *ImportFacade) GetErrors() []errordata.ErrorRecord {
	return i.errorData
}
