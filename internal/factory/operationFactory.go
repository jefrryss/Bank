package factory

import (
	"errors"
	"strings"
	"task2/domain/entities"
	"time"
)

type OperationFactory struct{}

func (f *OperationFactory) CreateOperation(
	id string,
	opType string,
	account *entities.BankAccount,
	category *entities.Category,
	amount float64,
	date time.Time,
	description string,
) (entities.Operation, error) {

	if strings.TrimSpace(id) == "" {
		return entities.Operation{}, errors.New("id не может быть пустым")
	}

	typeCatLower := strings.ToLower(strings.TrimSpace(opType))
	if typeCatLower != "доход" && typeCatLower != "расход" {
		return entities.Operation{}, errors.New("некорректный type категории")
	}

	if account == nil || strings.TrimSpace(account.ID) == "" {
		return entities.Operation{}, errors.New("аккаунт не может быть nil и должен иметь id")
	}

	if category == nil && strings.TrimSpace(category.ID) == "" {
		return entities.Operation{}, errors.New("категория не может быть nil и должна иметь id")
	}

	if amount <= 0 {
		return entities.Operation{}, errors.New("сумма должна быть положительной")
	}

	if date.IsZero() {
		return entities.Operation{}, errors.New("дата не может быть пустой")
	}

	return entities.Operation{
		ID:            strings.TrimSpace(id),
		TypeOperation: strings.ToLower(strings.TrimSpace(opType)),
		Account:       account,
		CategoryID:    category,
		Amount:        amount,
		Date:          date,
		Description:   strings.TrimSpace(description),
	}, nil
}
