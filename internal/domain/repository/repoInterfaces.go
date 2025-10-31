package repository

import "task2/internal/domain/entities"

type RepositoryBankAccounts interface {
	Save(account *entities.BankAccount) error
	Delete(id string) error
	Find(id string) (*entities.BankAccount, error)
	GetAll() ([]*entities.BankAccount, error)
}
type RepositoryCategory interface {
	Save(category *entities.Category) error
	Delete(id string) error
	Find(id string) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
}
type RepositoryOperations interface {
	Save(operation *entities.Operation) error
	Delete(id string) error
	Find(id string) (*entities.Operation, error)
	GetAll() ([]*entities.Operation, error)
}
