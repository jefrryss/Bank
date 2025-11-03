package bankmanager

import (
	"fmt"
	"strings"
	"task2/domain/entities"
	"task2/domain/errordata"
	"task2/domain/repository"
	"task2/internal/repository/memory"
	"time"
)

type BankManager struct {
	accountsRepo   repository.RepositoryBankAccounts
	categoriesRepo repository.RepositoryCategory
	operationsRepo repository.RepositoryOperations
	errorData      []errordata.ErrorRecord
}

func NewBankManager() *BankManager {
	return &BankManager{
		accountsRepo:   memory.NewRepositoryBankAccountsMemory(),
		categoriesRepo: memory.NewRepositoryCategoryMemory(),
		operationsRepo: memory.NewRepositoryOperMemory(),
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

// Работа с аккаунтами
func (m *BankManager) FindAccount(id string) (*entities.BankAccount, error) {
	return m.accountsRepo.Find(id)
}

func (m *BankManager) DeleteAccount(id string) error {
	return m.accountsRepo.Delete(id)
}

func (m *BankManager) GetAllAccounts() ([]*entities.BankAccount, error) {
	return m.accountsRepo.GetAll()
}

func (m *BankManager) GetErrors() []errordata.ErrorRecord {
	return m.errorData
}

// Работа с категориями
func (m *BankManager) FindCategory(id string) (*entities.Category, error) {
	return m.categoriesRepo.Find(id)
}

func (m *BankManager) DeleteCategory(id string) error {
	return m.categoriesRepo.Delete(id)
}

func (m *BankManager) GetAllCategories() ([]*entities.Category, error) {
	return m.categoriesRepo.GetAll()
}

// работа с операциями
func (m *BankManager) FindOperation(id string) (*entities.Operation, error) {
	return m.operationsRepo.Find(id)
}

func (m *BankManager) DeleteOperation(id string) error {
	return m.operationsRepo.Delete(id)
}

func (m *BankManager) GetAllOperations() ([]*entities.Operation, error) {
	return m.operationsRepo.GetAll()
}

// Дополнительные методы
func (m *BankManager) UpdateBalance(accountID string, newBalance float64) error {
	account, err := m.accountsRepo.Find(accountID)
	if err != nil {
		return err
	}

	account.Balance = newBalance

	return m.accountsRepo.Save(account)
}

func (m *BankManager) UpdateAccountFields(id, name string, balance float64) error {
	acc, err := m.accountsRepo.Find(id)
	if err != nil {
		return fmt.Errorf("счёт id=%s не найден: %w", id, err)
	}
	acc.Name = strings.TrimSpace(name)
	acc.Balance = balance
	return m.accountsRepo.Save(acc)
}

func (m *BankManager) UpdateCategoryFields(id, name, typeCat string) error {
	cat, err := m.categoriesRepo.Find(id)
	if err != nil {
		return fmt.Errorf("категория id=%s не найдена: %w", id, err)
	}
	cat.Name = strings.TrimSpace(name)
	cat.TypeCategory = strings.ToLower(strings.TrimSpace(typeCat))
	return m.categoriesRepo.Save(cat)
}

func (m *BankManager) UpdateOperationFields(
	id, opType, accountID, categoryID string,
	amount float64,
	date time.Time,
	description string,
) error {
	op, err := m.operationsRepo.Find(id)
	if err != nil {
		return fmt.Errorf("операция id=%s не найдена: %w", id, err)
	}
	acc, err := m.accountsRepo.Find(accountID)
	if err != nil {
		return fmt.Errorf("счёт id=%s не найден: %w", accountID, err)
	}
	cat, err := m.categoriesRepo.Find(categoryID)
	if err != nil {
		return fmt.Errorf("категория id=%s не найдена: %w", categoryID, err)
	}

	op.TypeOperation = strings.ToLower(strings.TrimSpace(opType))
	op.Account = acc
	op.CategoryID = cat
	op.Amount = amount
	op.Date = date
	op.Description = strings.TrimSpace(description)

	return m.operationsRepo.Save(op)
}
