package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ImportJSON struct {
	filePath  string
	errorData []errordata.ErrorRecord
	raw       map[string][]map[string]interface{}
}

func NewJSONParser() DataImporter {
	return &ImportJSON{}
}

func (importer *ImportJSON) GetPath() error {
	p := tea.NewProgram(newPathInputModel())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("ошибка интерфейса ввода: %v", err)
	}

	m := finalModel.(pathInputModel)

	if m.cancelled {
		return errors.New("ввод отменён пользователем")
	}

	if m.done {
		importer.filePath = m.filePath
		return nil
	}

	return errors.New("неизвестная ошибка при вводе пути")
}

func (i *ImportJSON) ParseData() error {
	data, err := os.ReadFile(i.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &i.raw)
}

func (i *ImportJSON) ParseBankAccounts() ([]entities.BankAccount, error) {
	var result []entities.BankAccount

	for _, acc := range i.raw["bank_accounts"] {
		result = append(result, entities.BankAccount{
			ID:      acc["id"].(string),
			Name:    acc["name"].(string),
			Balance: acc["balance"].(float64),
		})
	}

	return result, nil
}

func (i *ImportJSON) ParseCategories() ([]entities.Category, error) {
	var result []entities.Category

	for _, cat := range i.raw["categories"] {
		result = append(result, entities.Category{
			ID:           cat["id"].(string),
			TypeCategory: cat["type"].(string),
			Name:         cat["name"].(string),
		})
	}

	return result, nil
}

func (i *ImportJSON) ParseOperations() ([]entities.Operation, error) {
	var result []entities.Operation

	for _, op := range i.raw["operations"] {
		dateStr := op["date"].(string)
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			lineBytes, _ := json.Marshal(op)
			i.errorData = append(i.errorData, errordata.ErrorRecord{
				Line: string(lineBytes),
				Err:  errors.New("некорректная дата"),
			})
			continue
		}

		result = append(result, entities.Operation{
			ID:            op["id"].(string),
			TypeOperation: op["type"].(string),
			Account:       &entities.BankAccount{ID: op["account_id"].(string)},
			Amount:        op["amount"].(float64),
			Date:          parsedDate,
			Description:   op["description"].(string),
			CategoryID:    &entities.Category{ID: op["category_id"].(string)},
		})
	}

	return result, nil
}

func (i *ImportJSON) GetErrorData() []errordata.ErrorRecord {
	return i.errorData
}
