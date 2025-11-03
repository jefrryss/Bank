package importer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jefrryss/Bank/domain/errordata"

	"github.com/jefrryss/Bank/domain/ports"

	"time"

	"github.com/jefrryss/Bank/domain/entities"
)

type ImportJSON struct {
	filePath  string
	errorData []errordata.ErrorRecord
	raw       map[string][]map[string]interface{}
}

func NewJSONParser() ports.DataImporter {
	return &ImportJSON{}
}

func (i *ImportJSON) SetFilePath(path string) error {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, `"'`)

	if path == "" {
		return errors.New("путь не может быть пустым")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("файл %s не существует", path)
		}
		return fmt.Errorf("ошибка проверки файла: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("путь %s является директорией, а не файлом", path)
	}

	ext := filepath.Ext(path)
	if ext != ".json" {
		return fmt.Errorf("неверный формат файла: %s, ожидается .json", ext)
	}

	if info.Size() == 0 {
		return fmt.Errorf("файл %s пустой", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("нет прав на чтение файла: %v", err)
	}
	defer file.Close()

	firstByte := make([]byte, 1)
	_, err = file.Read(firstByte)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}

	if firstByte[0] != '{' && firstByte[0] != '[' {
		return fmt.Errorf("файл не является валидным JSON (должен начинаться с { или [)")
	}

	i.filePath = path
	return nil
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
