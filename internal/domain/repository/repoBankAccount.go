package repository

import "task2/internal/domain/entities"

type RepositoryBankAccounts interface {
	Save(account *entities.BankAccount) error
	Delete(id string) error
	Find(id string) (*entities.BankAccount, error)
	GetAll() ([]*entities.BankAccount, error)
}
