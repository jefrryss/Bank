package balancemanager

type BalanceManager interface {
	Recalculate() error
}
