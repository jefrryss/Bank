package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"task2/domain/entities"
)

// реализация ExportVisitor ( ExportCSV - ConcreateVisitor)
type ExportCSV struct {
	filePath string
}

func NewExportCSV() ExporterVisitor {
	return &ExportCSV{}
}

func (e *ExportCSV) SetFilePath(path string) error {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, `"'`)

	if path == "" {
		return errors.New("путь/имя файла не может быть пустым")
	}
	if filepath.Ext(path) != ".csv" {
		path += ".csv"
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить рабочую директорию: %v", err)
	}
	exportDir := filepath.Join(cwd, "exportFiles")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %v", exportDir, err)
	}

	finalPath := filepath.Join(exportDir, filepath.Base(path))

	f, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("нет прав на запись в файл %s: %v", finalPath, err)
	}
	_ = f.Close()

	e.filePath = finalPath
	return nil
}

func (exporter *ExportCSV) ExportBankAccount(accounts []*entities.BankAccount) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}
	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".csv"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]

	newFilePath := base + "_BackAccounts" + ext
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Name", "Balance"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, acc := range accounts {
		record := []string{
			acc.ID,
			acc.Name,
			fmt.Sprintf("%.2f", acc.Balance),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи: %v", err)
		}
	}

	return nil
}

func (exporter *ExportCSV) ExportCategory(categories []*entities.Category) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}
	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".csv"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]

	newFilePath := base + "_Category" + ext
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Type", "Name"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, cat := range categories {
		record := []string{
			cat.ID,
			cat.TypeCategory,
			cat.Name,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи записи: %v", err)
		}
	}

	return nil
}

func (exporter *ExportCSV) ExportOperation(operations []*entities.Operation) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}
	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".csv"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]

	newFilePath := base + "_Operation" + ext
	file, err := os.Create(newFilePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"ID", "Type", "AccountID", "Amount", "Date", "Description", "CategoryID"}); err != nil {
		return fmt.Errorf("ошибка записи заголовка: %v", err)
	}

	for _, op := range operations {
		record := []string{
			op.ID,
			op.TypeOperation,
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
