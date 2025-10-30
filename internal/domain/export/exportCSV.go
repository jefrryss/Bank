package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"task2/internal/domain/entities"
)

type ExportCSV struct {
	filePath string
}

func (exporter *ExportCSV) GetPath() error {

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить рабочую директорию: %v", err)
	}
	folder := wd + "/../exportFiles"

	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку: %v", err)
	}

	var filename string
	for {
		fmt.Print("Введите название файла для экспорта (например bank_accounts.csv): ")
		fmt.Scanln(&filename)
		filename = strings.TrimSpace(filename)

		if filename == "" {
			fmt.Println("Название файла не может быть пустым")
			continue
		}

		fullPath := folder + "/" + filename

		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("Файл %s уже существует. Перезаписать? (y/n): ", fullPath)
			var answer string
			fmt.Scanln(&answer)
			answer = strings.ToLower(strings.TrimSpace(answer))
			if answer == "y" {
				exporter.filePath = fullPath
				break
			} else {
				fmt.Println("Введите новое имя файла")
				continue
			}
		} else {
			exporter.filePath = fullPath
			break
		}
	}

	return nil
}

func (exporter *ExportCSV) ExportBankAccount(accounts *[]entities.BankAccount) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	file, err := os.Create(exporter.filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Name", "Balance"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, acc := range *accounts {
		record := []string{
			acc.ID,
			acc.Name,
			fmt.Sprintf("%.2f", acc.Balance()),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи: %v", err)
		}
	}

	return nil
}

func (exporter *ExportCSV) ExportCategory(categories *[]entities.Category) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	file, err := os.Create(exporter.filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Type", "Name"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, cat := range *categories {
		record := []string{
			cat.ID,
			cat.TypeCategory(),
			cat.Name,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи записи: %v", err)
		}
	}

	return nil
}

func (exporter *ExportCSV) ExportOperation(operations *[]entities.Operation) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	file, err := os.Create(exporter.filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Type", "AccountID", "Amount", "Date", "Description", "CategoryID"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, op := range *operations {
		record := []string{
			op.ID,
			op.TypeOperation(),
			op.Account.ID,
			fmt.Sprintf("%.2f", op.Amount),
			op.Date.Format("2006-01-02"),
			op.Description,
			"",
		}
		if op.CategoryID != nil {
			record[6] = op.CategoryID.ID
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи записи: %v", err)
		}
	}

	return nil
}
