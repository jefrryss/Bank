package bankmanager

import (
	"fmt"
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
	"task2/internal/domain/repository"
)

type BankManager struct {
	accountsRepo   repository.RepositoryBankAccounts
	categoriesRepo repository.RepositoryCategory
	operationsRepo repository.RepositoryOperations
	errorData      []errordata.ErrorRecord
}

func NewBankManager() *BankManager {
	return &BankManager{
		accountsRepo:   repository.NewRepositoryBankAccountsMemory(),
		categoriesRepo: repository.NewRepositoryCategoryMemory(),
		operationsRepo: repository.NewRepositoryOperMemory(),
	}
}

// Добавление нескольких аккаунтов с проверкой уникальности
func (m *BankManager) AddAccounts(accounts []*entities.BankAccount) {
	for _, acc := range accounts {
		if _, err := m.accountsRepo.Find(acc.ID); err == nil {
			m.errorData = append(m.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", acc),
				Err:  fmt.Errorf("аккаунт с id %s уже существует", acc.ID),
			})
			continue
		}
		_ = m.accountsRepo.Save(acc)
	}
}

// Добавление нескольких категорий с проверкой уникальности
func (m *BankManager) AddCategories(categories []*entities.Category) {
	for _, cat := range categories {
		if _, err := m.categoriesRepo.Find(cat.ID); err == nil {
			m.errorData = append(m.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", cat),
				Err:  fmt.Errorf("категория с id %s уже существует", cat.ID),
			})
			continue
		}
		_ = m.categoriesRepo.Save(cat)
	}
}

// Добавление нескольких операций с проверкой уникальности
func (m *BankManager) AddOperations(operations []*entities.Operation) {
	for _, op := range operations {
		if _, err := m.operationsRepo.Find(op.ID); err == nil {
			m.errorData = append(m.errorData, errordata.ErrorRecord{
				Line: fmt.Sprintf("%v", op),
				Err:  fmt.Errorf("операция с id %s уже существует", op.ID),
			})
			continue
		}
		_ = m.operationsRepo.Save(op)
	}
}
func (m *BankManager) FindAccount(id string) (*entities.BankAccount, error) {
	return m.accountsRepo.Find(id)
}

func (m *BankManager) DeleteAccount(id string) error {
	return m.accountsRepo.Delete(id)
}

func (m *BankManager) GetAllAccounts() ([]*entities.BankAccount, error) {
	return m.accountsRepo.GetAll()
}

// Получить все ошибки
func (m *BankManager) GetErrors() []errordata.ErrorRecord {
	return m.errorData
}

func (m *BankManager) FindCategory(id string) (*entities.Category, error) {
	return m.categoriesRepo.Find(id)
}

func (m *BankManager) DeleteCategory(id string) error {
	return m.categoriesRepo.Delete(id)
}

func (m *BankManager) GetAllCategories() ([]*entities.Category, error) {
	return m.categoriesRepo.GetAll()
}
func (m *BankManager) FindOperation(id string) (*entities.Operation, error) {
	return m.operationsRepo.Find(id)
}

func (m *BankManager) DeleteOperation(id string) error {
	return m.operationsRepo.Delete(id)
}

func (m *BankManager) GetAllOperations() ([]*entities.Operation, error) {
	return m.operationsRepo.GetAll()
}
