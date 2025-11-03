package analitic

import (
	"fmt"
	"strings"

	"github.com/jefrryss/Bank/domain/entities"
)

type GroupByCategoryCommand struct {
	Operations []*entities.Operation
	Categories []*entities.Category
	Result     string
}

func NewGroupByCategoryCommand(op []*entities.Operation, cat []*entities.Category) Command {
	return &GroupByCategoryCommand{
		Operations: op,
		Categories: cat,
	}
}
func (c *GroupByCategoryCommand) Execute() (string, error) {
	categoryMap := make(map[string]string)
	for _, cat := range c.Categories {
		categoryMap[cat.ID] = cat.Name
	}

	incomeByCategory := make(map[string]float64)
	expenseByCategory := make(map[string]float64)

	for _, op := range c.Operations {
		categoryID := ""
		if op.CategoryID != nil {
			categoryID = op.CategoryID.ID
		}

		catName := categoryMap[categoryID]
		if catName == "" {
			catName = "Без категории"
		}

		if op.TypeOperation == "доход" {
			incomeByCategory[catName] += op.Amount
		} else {
			expenseByCategory[catName] += op.Amount
		}
	}

	var result strings.Builder
	result.WriteString("=== ДОХОДЫ ПО КАТЕГОРИЯМ ===\n")
	if len(incomeByCategory) == 0 {
		result.WriteString("Нет доходов\n")
	} else {
		for cat, amount := range incomeByCategory {
			result.WriteString(fmt.Sprintf("%s: %.2f\n", cat, amount))
		}
	}

	result.WriteString("\n=== РАСХОДЫ ПО КАТЕГОРИЯМ ===\n")
	if len(expenseByCategory) == 0 {
		result.WriteString("Нет расходов\n")
	} else {
		for cat, amount := range expenseByCategory {
			result.WriteString(fmt.Sprintf("%s: %.2f\n", cat, amount))
		}
	}

	return result.String(), nil
}
