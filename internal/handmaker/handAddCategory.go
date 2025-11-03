package handmaker

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jefrryss/Bank/domain/entities"
	"github.com/jefrryss/Bank/internal/bankmanager"
	"github.com/jefrryss/Bank/internal/factory"
)

type Controller struct {
	mgr     *bankmanager.BankManager
	accFact *factory.BankAccount
	catFact *factory.CategoryFactory
	opFact  *factory.OperationFactory
}

func NewController(mgr *bankmanager.BankManager) *Controller {
	return &Controller{
		mgr:     mgr,
		accFact: &factory.BankAccount{},
		catFact: &factory.CategoryFactory{},
		opFact:  &factory.OperationFactory{},
	}
}

func (c *Controller) CreateAccount(id, name, balanceStr string) error {
	if _, err := c.mgr.FindAccount(id); err == nil {
		return fmt.Errorf("аккаунт с id %q уже существует", id)
	}

	bal, err := parseFloat(balanceStr)
	if err != nil {
		return fmt.Errorf("balance: %w", err)
	}
	acc, err := c.accFact.CreateBankAccount(id, name, bal)
	if err != nil {
		return err
	}

	c.mgr.AddAccounts([]*entities.BankAccount{&acc})
	return nil
}

func (c *Controller) UpdateAccount(id, name, balanceStr string) error {
	bal, err := parseFloat(balanceStr)
	if err != nil {
		return fmt.Errorf("balance: %w", err)
	}
	if _, err := c.accFact.CreateBankAccount(id, name, bal); err != nil {
		return err
	}
	return c.mgr.UpdateAccountFields(id, name, bal)
}

func (c *Controller) DeleteAccount(id string) error {

	if _, err := c.mgr.FindAccount(id); err != nil {
		return fmt.Errorf("аккаунт с id %q не найден", id)
	}
	return c.mgr.DeleteAccount(id)
}

func (c *Controller) CreateCategory(id, name, typeCat string) error {
	if _, err := c.mgr.FindCategory(id); err == nil {
		return fmt.Errorf("категория с id %q уже существует", id)
	}

	cat, err := c.catFact.CreateCategory(id, name, typeCat)
	if err != nil {
		return err
	}
	c.mgr.AddCategories([]*entities.Category{&cat})
	return nil
}

func (c *Controller) UpdateCategory(id, name, typeCat string) error {
	if _, err := c.catFact.CreateCategory(id, name, typeCat); err != nil {
		return err
	}
	return c.mgr.UpdateCategoryFields(id, name, typeCat)
}

func (c *Controller) DeleteCategory(id string) error {

	if _, err := c.mgr.FindCategory(id); err != nil {
		return fmt.Errorf("категория с id %q не найдена", id)
	}
	return c.mgr.DeleteCategory(id)
}

func (c *Controller) CreateOperation(
	id, opType, accountID, categoryID, amountStr, dateStr, description string,
) error {
	if _, err := c.mgr.FindOperation(id); err == nil {
		return fmt.Errorf("операция с id %q уже существует", id)
	}

	acc, err := c.mgr.FindAccount(accountID)
	if err != nil {
		return fmt.Errorf("account id=%q не найден: %w", accountID, err)
	}
	cat, err := c.mgr.FindCategory(categoryID)
	if err != nil {
		return fmt.Errorf("category id=%q не найдена: %w", categoryID, err)
	}

	amt, err := parseFloat(amountStr)
	if err != nil {
		return fmt.Errorf("amount: %w", err)
	}
	dt, err := parseDate(dateStr)
	if err != nil {
		return fmt.Errorf("date: %w", err)
	}
	op, err := c.opFact.CreateOperation(id, opType, acc, cat, amt, dt, description)
	if err != nil {
		return err
	}

	c.mgr.AddOperations([]*entities.Operation{&op})
	return nil
}

func (c *Controller) UpdateOperation(
	id, opType, accountID, categoryID, amountStr, dateStr, description string,
) error {
	amt, err := parseFloat(amountStr)
	if err != nil {
		return fmt.Errorf("amount: %w", err)
	}
	dt, err := parseDate(dateStr)
	if err != nil {
		return fmt.Errorf("date: %w", err)
	}

	acc, err := c.mgr.FindAccount(accountID)
	if err != nil {
		return fmt.Errorf("account id=%q: %w", accountID, err)
	}
	cat, err := c.mgr.FindCategory(categoryID)
	if err != nil {
		return fmt.Errorf("category id=%q: %w", categoryID, err)
	}

	if _, err := c.opFact.CreateOperation(id, opType, acc, cat, amt, dt, description); err != nil {
		return err
	}

	return c.mgr.UpdateOperationFields(id, opType, accountID, categoryID, amt, dt, description)
}

func (c *Controller) DeleteOperation(id string) error {

	if _, err := c.mgr.FindOperation(id); err != nil {
		return fmt.Errorf("операция с id %q не найдена", id)
	}
	return c.mgr.DeleteOperation(id)
}

func parseFloat(s string) (float64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", ".")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("некорректное число %q", s)
	}
	return val, nil
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("ожидается YYYY-MM-DD (%q)", s)
	}
	return d, nil
}
