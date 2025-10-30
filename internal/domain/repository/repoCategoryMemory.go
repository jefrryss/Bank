package repository

import (
	"errors"
	"task2/internal/domain/entities"
)

//Реалиализация интерфейса RepositoryCategory

type RepositoryCategoryMemory struct {
	data map[string]*entities.Category
}

func NewRepositoryCategoryMemory() RepositoryCategory {
	return &RepositoryCategoryMemory{data: make(map[string]*entities.Category)}
}

func (repo *RepositoryCategoryMemory) Save(category *entities.Category) error {
	repo.data[category.ID] = category
	return nil
}

func (repo *RepositoryCategoryMemory) Delete(id string) error {
	delete(repo.data, id)
	return nil
}

func (repo *RepositoryCategoryMemory) Find(id string) (*entities.Category, error) {
	if item, checker := repo.data[id]; checker {
		return item, nil
	}
	return nil, errors.New("элемент не найден")
}

func (repo *RepositoryCategoryMemory) GetAll() ([]*entities.Category, error) {
	var result []*entities.Category
	for _, item := range repo.data {
		result = append(result, item)
	}
	return result, nil
}
