package importer

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task2/internal/domain/entities"
	"task2/internal/domain/errordata"
	"time"
)

type ImportCSV struct {
	filePath  string
	records   [][]string
	errorData []errordata.ErrorRecord
}

func NewCSVParser() DataImporter {
	return &ImportCSV{}
}

func (importer *ImportCSV) GetPath() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите путь к CSV файлу: ")

	filePath, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	filePath = strings.TrimSpace(filePath)

	if filePath == "" {
		return "", errors.New("путь не может быть пустым")
	}

	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("файл или каталог не существует: %s", filePath)
		}
	}

	importer.filePath = filePath
	return filePath, nil
}

// ParseData читает данные из файла и сохраняет в структуру (разделитель ",")
func (importer *ImportCSV) ParseData() error {
	if importer.filePath == "" {
		return errors.New("путь к файлу не установлен")
	}

	file, err := os.Open(importer.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("пустой файл")
	}

	importer.records = records
	return nil
}

func (importer *ImportCSV) ParseBankAccounts() ([]entities.BankAccount, error) {
	var accounts []entities.BankAccount

	for i := 1; i < len(importer.records); i++ {
		record := importer.records[i]

		for len(record) < 10 {
			record = append(record, "")
		}

		if strings.TrimSpace(record[0]) != "bank_account" {
			continue
		}

		idInt, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в id"),
			})
			continue
		}

		balance, err := parseFloat(record[9])
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в balance"),
			})
			continue
		}

		account := entities.BankAccount{
			ID:   strconv.Itoa(idInt),
			Name: strings.TrimSpace(record[8]),
		}

		if err := account.SetBalance(balance); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  err,
			})
			continue
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// ParseCategories парсит категории
func (importer *ImportCSV) ParseCategories() ([]entities.Category, error) {
	var categories []entities.Category

	for i := 1; i < len(importer.records); i++ {
		record := importer.records[i]

		for len(record) < 10 {
			record = append(record, "")
		}

		if strings.TrimSpace(record[0]) != "category" {
			continue
		}

		idInt, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в id"),
			})
			continue
		}

		category := entities.Category{
			ID:   strconv.Itoa(idInt),
			Name: strings.TrimSpace(record[8]),
		}

		if err := category.SetTypeCategory(strings.TrimSpace(record[2])); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  err,
			})
			continue
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// ParseOperations парсит операции
func (importer *ImportCSV) ParseOperations() ([]entities.Operation, error) {
	var operations []entities.Operation

	for i := 1; i < len(importer.records); i++ {
		record := importer.records[i]

		for len(record) < 10 {
			record = append(record, "")
		}

		if strings.TrimSpace(record[0]) != "operation" {
			continue
		}

		idInt, err := strconv.Atoi(strings.TrimSpace(record[1]))
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в id"),
			})
			continue
		}

		bankAccountID, err := strconv.Atoi(strings.TrimSpace(record[3]))
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в BankAccountID"),
			})
			continue
		}

		amount, err := parseFloat(record[4])
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в amount"),
			})
			continue
		}

		date, err := time.Parse("2006-01-02", strings.TrimSpace(record[5]))
		if err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  errors.New("ошибка в Дате"),
			})
			continue
		}

		categoryID := ""
		if strings.TrimSpace(record[7]) != "" {
			catID, err := strconv.Atoi(strings.TrimSpace(record[7]))
			if err != nil {
				importer.errorData = append(importer.errorData, errordata.ErrorRecord{
					Line: record,
					Err:  errors.New("ошибка в Category id"),
				})
				continue
			}
			categoryID = strconv.Itoa(catID)
		}

		operation := entities.Operation{
			ID:          strconv.Itoa(idInt),
			Account:     &entities.BankAccount{ID: strconv.Itoa(bankAccountID)},
			Amount:      amount,
			Date:        date,
			Description: strings.TrimSpace(record[6]),
			CategoryID:  &entities.Category{ID: categoryID},
		}
		if err := operation.SetTypeOperation(strings.TrimSpace(record[2])); err != nil {
			importer.errorData = append(importer.errorData, errordata.ErrorRecord{
				Line: record,
				Err:  err,
			})
			continue
		}
		operations = append(operations, operation)
	}

	return operations, nil
}
func (importer *ImportCSV) GetErrorData() []errordata.ErrorRecord { return importer.errorData }

// parseFloat парсит строку в float64
func parseFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}
