package repository

import "github.com/jefrryss/Bank/domain/entities"

type RepositoryBankAccounts interface {
	Save(account *entities.BankAccount) error
	Delete(id string) error
	Find(id string) (*entities.BankAccount, error)
	GetAll() ([]*entities.BankAccount, error)
}
