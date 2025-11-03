package balancemanager

import "task2/internal/bankmanager"

type ManualBalanceManager struct {
	Bank       *bankmanager.BankManager
	AccountID  string
	NewBalance float64
}

func NewManualBalanceManager(bank *bankmanager.BankManager, accountID string, newBalance float64) *ManualBalanceManager {
	return &ManualBalanceManager{
		Bank:       bank,
		AccountID:  accountID,
		NewBalance: newBalance,
	}
}

func (m *ManualBalanceManager) Recalculate() error {
	return m.Bank.UpdateBalance(m.AccountID, m.NewBalance)
}
