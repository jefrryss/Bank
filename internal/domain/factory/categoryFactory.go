package factory

import (
	"errors"
	"strings"
	"task2/internal/domain/entities"
)

type CategoryFactory struct {
}

func (c *CategoryFactory) CreateCategory(id string, name string, typeCat string) (entities.Category, error) {
	if strings.TrimSpace(id) == "" {
		return entities.Category{}, errors.New("id не может быть пустым")
	}
	if strings.TrimSpace(name) == "" {
		return entities.Category{}, errors.New("name не может быть пустым")
	}
	typeCatLower := strings.ToLower(strings.TrimSpace(typeCat))
	if typeCatLower != "доход" && typeCatLower != "расход" {
		return entities.Category{}, errors.New("некорректный type категории")
	}
	return entities.Category{
		ID:           strings.TrimSpace(id),
		Name:         strings.TrimSpace(name),
		TypeCategory: strings.ToLower(strings.TrimSpace(typeCat)),
	}, nil
}
