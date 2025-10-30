package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
	"time"
)

type jsonData struct {
	BankAccounts []struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	} `json:"bank_accounts"`
	Categories []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"categories"`
	Operations []struct {
		ID          string  `json:"id"`
		Type        string  `json:"type"`
		AccountID   string  `json:"account_id"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Description string  `json:"description"`
		CategoryID  string  `json:"category_id"`
	} `json:"operations"`
}

type ImportJSON struct {
	filePath    string
	errorData   []errordata.ErrorRecord
	recordsJSON jsonData
}

func NewJSONParser() DataImporter {
	return &ImportJSON{}
}

func (importer *ImportJSON) GetPath() (string, error) {
	fmt.Print("Введите путь к JSON файлу: ")
	var filePath string
	fmt.Scanln(&filePath)

	if filePath == "" {
		return "", errors.New("путь не может быть пустым")
	}

	importer.filePath = strings.TrimSpace(filePath)
	return importer.filePath, nil
}

func (importer *ImportJSON) ParseData() error {
	if importer.filePath == "" {
		return errors.New("путь к файлу не установлен")
	}

	dataBytes, err := os.ReadFile(importer.filePath)
	if err != nil {
		return err
	}

	var jd jsonData
	if err := json.Unmarshal(dataBytes, &jd); err != nil {
		return err
	}

	importer.recordsJSON = jd
	return nil
}

func (importer *ImportJSON) ParseBankAccounts() ([]entities.BankAccount, error) {
	var accounts []entities.BankAccount
	for _, a := range importer.recordsJSON.BankAccounts {
		account := entities.BankAccount{
			ID:   a.ID,
			Name: a.Name,
		}
		if err := account.SetBalance(a.Balance); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: []string{a.ID, a.Name, fmt.Sprintf("%f", a.Balance)},
				Err:  err,
			})
			continue
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (importer *ImportJSON) ParseCategories() ([]entities.Category, error) {
	var categories []entities.Category
	for _, c := range importer.recordsJSON.Categories {
		category := entities.Category{
			ID:   c.ID,
			Name: c.Name,
		}
		if err := category.SetTypeCategory(c.Type); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: []string{c.ID, c.Type, c.Name},
				Err:  err,
			})
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (importer *ImportJSON) ParseOperations() ([]entities.Operation, error) {
	var operations []entities.Operation
	for _, o := range importer.recordsJSON.Operations {
		date, err := time.Parse("2006-01-02", o.Date)
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: []string{o.ID, o.Type, o.AccountID, o.Date},
				Err:  errors.New("ошибка в дате"),
			})
			continue
		}

		op := entities.Operation{
			ID:          o.ID,
			Account:     &entities.BankAccount{ID: o.AccountID},
			Amount:      o.Amount,
			Date:        date,
			Description: o.Description,
			CategoryID:  &entities.Category{ID: o.CategoryID},
		}

		if err := op.SetTypeOperation(o.Type); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: []string{o.ID, o.Type},
				Err:  err,
			})
			continue
		}

		operations = append(operations, op)
	}
	return operations, nil
}

func (importer *ImportJSON) GetErrorData() []errordata.ErrorRecord {
	return importer.errorData
}
