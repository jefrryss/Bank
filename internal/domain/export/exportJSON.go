package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"task2/internal/domain/entities"
)

type ExportJSON struct {
	filePath string
}

func NewExportJSON() ExporterVisitor {
	return &ExportJSON{}
}
func (exporter *ExportJSON) GetPath() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("не удалось получить рабочую директорию: %v", err)
	}

	folder := filepath.Join(wd, "..", "exportFiles")
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку: %v", err)
	}

	var filename string
	for {
		fmt.Print("Введите название файла для экспорта (например data.json): ")
		fmt.Scanln(&filename)
		filename = strings.TrimSpace(filename)

		if filename == "" {
			fmt.Println("Название файла не может быть пустым")
			continue
		}

		if filepath.Ext(filename) == "" {
			filename += ".json"
		}

		fullPath := filepath.Join(folder, filename)

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

func (exporter *ExportJSON) ExportBankAccount(accounts *[]*entities.BankAccount) error {
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

	exportData := make([]bankAccountExport, 0, len(*accounts))
	for _, acc := range *accounts {
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

func (exporter *ExportJSON) ExportCategory(categories *[]*entities.Category) error {
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

	exportData := make([]categoryExport, 0, len(*categories))
	for _, cat := range *categories {
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

func (exporter *ExportJSON) ExportOperation(operations *[]*entities.Operation) error {
	if exporter.filePath == "" {
		return fmt.Errorf("путь к файлу не задан")
	}

	ext := filepath.Ext(exporter.filePath)
	if ext == "" {
		ext = ".json"
	}
	base := exporter.filePath[:len(exporter.filePath)-len(ext)]
	newFilePath := base + "_Operation.json"

	// Тип для экспорта с вложенными объектами и нужными полями
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

	exportData := make([]operationExport, 0, len(*operations))
	for _, op := range *operations {
		acc := accountExport{
			ID:      op.Account.ID,
			Name:    op.Account.Name,
			Balance: op.Account.Balance,
		}

		cat := categoryExport{
			ID:   "",
			Type: "",
			Name: "",
		}
		if op.CategoryID != nil {
			cat.ID = op.CategoryID.ID
			cat.Type = op.CategoryID.TypeCategory
			cat.Name = op.CategoryID.Name
		}

		exportData = append(exportData, operationExport{
			ID:          op.ID,
			Type:        op.TypeOperation,
			Account:     acc,
			Amount:      op.Amount,
			Date:        op.Date.Format("2006-01-02"),
			Description: op.Description,
			Category:    cat,
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
