package main

import (
	"fmt"
	"log"
	"task2/internal/domain/entities"
	"task2/internal/domain/export" // ваш пакет с ExportCSV
)

func main() {
	// Создаем тестовые банковские счета
	accounts := []entities.BankAccount{
		{ID: "1", Name: "Основной счет"},
		{ID: "2", Name: "Кредитная карта"},
		{ID: "3", Name: "Накопительный счет"},
	}

	// Устанавливаем баланс
	for i := range accounts {
		var err error
		switch accounts[i].ID {
		case "1":
			err = accounts[i].SetBalance(25000.0)
		case "2":
			err = accounts[i].SetBalance(-5000.0) // ошибка, покажет проверку
		case "3":
			err = accounts[i].SetBalance(100000.0)
		}
		if err != nil {
			fmt.Printf("Ошибка установки баланса для счета %s: %v\n", accounts[i].ID, err)
		}
	}

	// Создаем Exporter
	exporter := export.ExportCSV{}

	// Получаем путь для файла
	if err := exporter.GetPath(); err != nil {
		log.Fatalf("Ошибка с путем к файлу: %v", err)
	}

	// Экспортируем BankAccount
	if err := exporter.ExportBankAccount(&accounts); err != nil {
		log.Fatalf("Ошибка экспорта: %v", err)
	}

	fmt.Println("Экспорт завершен успешно!")
}
