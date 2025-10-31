package factory

import (
	"errors"
	"strings"
	"task2/internal/domain/entities"
)

type BankAccount struct {
}

func (b *BankAccount) CreateBankAccount(id string, name string, balance float64) (entities.BankAccount, error) {

	if strings.TrimSpace(id) == "" {
		return entities.BankAccount{}, errors.New("id не может быть пустым")
	}

	if strings.TrimSpace(name) == "" {
		return entities.BankAccount{}, errors.New("имя не может быть пустым")
	}

	if balance < 0 {
		return entities.BankAccount{}, errors.New("баланс не может быть отрицательным")
	}
	return entities.BankAccount{
		ID:      id,
		Name:    name,
		Balance: balance,
	}, nil
}
