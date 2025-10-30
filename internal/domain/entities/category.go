package entities

import (
	"errors"
	"strings"
)

type Category struct {
	ID           string
	typeCategory string
	Name         string
}

func (c *Category) SetTypeCategory(typeInput string) error {
	typeInput = strings.ToLower(strings.TrimSpace(typeInput))

	switch typeInput {
	case "доход":
		c.typeCategory = "доход"
		return nil
	case "расход":
		c.typeCategory = "расход"
		return nil
	default:
		return errors.New("категория может быть только 'Доход' или 'Расход'")
	}
}

func (c *Category) TypeCategory() string { return c.typeCategory }
