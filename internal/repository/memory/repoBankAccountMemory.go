package memory

import (
	"errors"

	"github.com/jefrryss/Bank/domain/entities"
	"github.com/jefrryss/Bank/domain/repository"
)

//Реалиализация интерфейса RepositoryBankAccounts

type RepositoryBankAccountsMemory struct {
	data map[string]*entities.BankAccount
}

func NewRepositoryBankAccountsMemory() repository.RepositoryBankAccounts {
	return &RepositoryBankAccountsMemory{data: make(map[string]*entities.BankAccount)}
}

func (repo *RepositoryBankAccountsMemory) Save(account *entities.BankAccount) error {
	repo.data[account.ID] = account
	return nil
}

func (repo *RepositoryBankAccountsMemory) Delete(id string) error {
	delete(repo.data, id)
	return nil
}

func (repo *RepositoryBankAccountsMemory) Find(id string) (*entities.BankAccount, error) {
	if item, checker := repo.data[id]; checker {
		return item, nil
	}
	return nil, errors.New("элемент не найден")
}

func (repo *RepositoryBankAccountsMemory) GetAll() ([]*entities.BankAccount, error) {
	var result []*entities.BankAccount
	for _, item := range repo.data {
		result = append(result, item)
	}
	return result, nil
}
