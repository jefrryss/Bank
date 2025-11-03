package balancemanager

import "task2/internal/bankmanager"

type AutoBalanceManager struct {
	Bank *bankmanager.BankManager
}

func NewAutoBalanceManager(bank *bankmanager.BankManager) *AutoBalanceManager {
	return &AutoBalanceManager{Bank: bank}
}

func (a *AutoBalanceManager) Recalculate() error {
	accounts, err := a.Bank.GetAllAccounts()
	if err != nil {
		return err
	}

	operations, err := a.Bank.GetAllOperations()
	if err != nil {
		return err
	}

	for _, acc := range accounts {
		var balance float64

		for _, op := range operations {
			if op.Account.ID != acc.ID {
				continue
			}

			switch op.TypeOperation {
			case "доход":
				balance += op.Amount
			case "расход":
				balance -= op.Amount
			}
		}

		if err := a.Bank.UpdateBalance(acc.ID, balance); err != nil {
			return err
		}
	}

	return nil
}
