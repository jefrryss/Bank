package export

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jefrryss/Bank/domain/entities"
)

// реализация ExportVisitor ( ExportJSON - ConcreateVisitor)
type ExportJSON struct {
	filePath string
}

func NewExportJSON() ExporterVisitor {
	return &ExportJSON{}
}

func (e *ExportJSON) SetFilePath(path string) error {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, `"'`)

	if path == "" {
		return errors.New("путь/имя файла не может быть пустым")
	}
	if filepath.Ext(path) != ".json" {
		path += ".json"
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить рабочую директорию: %v", err)
	}
	exportDir := filepath.Join(cwd, "exportFiles")
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %v", exportDir, err)
	}

	finalPath := filepath.Join(exportDir, filepath.Base(path))

	testFile, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("нет прав на запись в файл %s: %v", finalPath, err)
	}
	testFile.Close()

	e.filePath = finalPath
	return nil
}

func (exporter *ExportJSON) ExportBankAccount(accounts []*entities.BankAccount) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".json"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]
	newFilePath := base + "_BackAccounts.json"

	type bankAccountExport struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}

	exportData := make([]bankAccountExport, 0, len(accounts))
	for _, acc := range accounts {
		exportData = append(exportData, bankAccountExport{
			ID:      acc.ID,
			Name:    acc.Name,
			Balance: acc.Balance,
		})
	}

	dataBytes, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	if err := os.WriteFile(newFilePath, dataBytes, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

func (exporter *ExportJSON) ExportCategory(categories []*entities.Category) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".json"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]
	newFilePath := base + "_Category.json"

	type categoryExport struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}

	exportData := make([]categoryExport, 0, len(categories))
	for _, cat := range categories {
		exportData = append(exportData, categoryExport{
			ID:   cat.ID,
			Type: cat.TypeCategory,
			Name: cat.Name,
		})
	}

	dataBytes, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	if err := os.WriteFile(newFilePath, dataBytes, 0644); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	return nil
}

func (exporter *ExportJSON) ExportOperation(operations []*entities.Operation) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	var ops []*entities.Operation
	if operations != nil {
		ops = operations
	} else {
		ops = make([]*entities.Operation, 0)
	}

	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".json"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]
	newFilePath := base + "_Operation.json"

	type accountExport struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Balance float64 `json:"balance"`
	}
	type categoryExport struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}
	type operationExport struct {
		ID          string         `json:"id"`
		Type        string         `json:"type"`
		Account     accountExport  `json:"account"`
		Amount      float64        `json:"amount"`
		Date        string         `json:"date"`
		Description string         `json:"description"`
		Category    categoryExport `json:"category"`
	}

	exportData := make([]operationExport, 0, len(ops))
	for i := range ops {
		op := ops[i]

		acc := accountExport{}
		if op.Account != nil {
			acc.ID = op.Account.ID
			acc.Name = op.Account.Name
			acc.Balance = op.Account.Balance
		}

		cat := categoryExport{}
		if op.CategoryID != nil {
			cat.ID = op.CategoryID.ID
			cat.Type = op.CategoryID.TypeCategory
			cat.Name = op.CategoryID.Name
		}

		dateStr := ""
		if !op.Date.IsZero() {
			dateStr = op.Date.Format("2006-01-02")
		}

		exportData = append(exportData, operationExport{
			ID:          op.ID,
			Type:        op.TypeOperation,
			Account:     acc,
			Amount:      op.Amount,
			Date:        dateStr,
			Description: op.Description,
			Category:    cat,
		})
	}

	dataBytes, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}
	if err := os.WriteFile(newFilePath, dataBytes, 0o644); err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}
	return nil
}
