package balancemanager

import (
	"task2/internal/domain/repository"
)

type BalanceAdjuster interface {
	Recalculate() error
}

// Автоматический пересчёт баланса по всем операциям
type AutoBalanceAdjuster struct {
	AccountRepo   repository.RepositoryBankAccounts
	OperationRepo repository.RepositoryOperations
}

func (a *AutoBalanceAdjuster) Recalculate() error {
	accounts, err := a.AccountRepo.GetAll()
	if err != nil {
		return err
	}

	ops, err := a.OperationRepo.GetAll()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		var balance float64
		for _, op := range ops {
			if op.Account.ID != account.ID {
				continue
			}

			if op.TypeOperation == "доход" {
				balance += op.Amount
			} else if op.TypeOperation == "расход" {
				balance -= op.Amount
			}
		}

		account.Balance = balance
	}
	return nil
}

// Ручная корректировка баланса конкретного аккаунта
type ManualBalanceAdjuster struct {
	AccountRepo repository.RepositoryBankAccounts
	AccountID   string
	NewBalance  float64
}

func (m *ManualBalanceAdjuster) Recalculate() error {
	account, err := m.AccountRepo.Find(m.AccountID)
	if err != nil {
		return err
	}

	account.Balance = m.NewBalance
	return nil
}
