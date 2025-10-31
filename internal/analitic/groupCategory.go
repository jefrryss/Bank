package analitic

import (
	"fmt"
	"task2/internal/domain/entities"
)

// Команда для группировки операций по категориям
type GroupByCategoryCommand struct {
	Operations []*entities.Operation
	Categories []*entities.Category
}

// Конструктор команды, принимает слайсы данных напрямую
func NewGroupByCategoryCommand(ops []*entities.Operation, cats []*entities.Category) *GroupByCategoryCommand {
	return &GroupByCategoryCommand{
		Operations: ops,
		Categories: cats,
	}
}

// Выполнение команды
func (c *GroupByCategoryCommand) Execute() error {
	result := make(map[string]float64)

	// Создаем карту id -> имя категории
	catMap := make(map[string]string)
	for _, category := range c.Categories {
		catMap[category.ID] = category.Name
	}

	// Суммируем операции по категориям
	for _, op := range c.Operations {
		if op.CategoryID == nil {
			continue
		}
		name := catMap[op.CategoryID.ID]
		result[name] += op.Amount
	}

	// Вывод результата
	fmt.Println("Группировка по категориям:")
	for name, sum := range result {
		fmt.Printf("%s: %.2f\n", name, sum)
	}

	return nil
}
