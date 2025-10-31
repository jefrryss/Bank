package repository

import (
	"errors"
	"task2/internal/domain/entities"
)

//Реалиализация интерфейса RepositoryOperations

type RepositoryOperMemory struct {
	data map[string]*entities.Operation
}

func NewRepositoryOperMemory() RepositoryOperations {
	return &RepositoryOperMemory{data: make(map[string]*entities.Operation)}
}

func (repo *RepositoryOperMemory) Save(operation *entities.Operation) error {
	repo.data[operation.ID] = operation
	return nil
}

func (repo *RepositoryOperMemory) Delete(id string) error {
	delete(repo.data, id)
	return nil
}

func (repo *RepositoryOperMemory) Find(id string) (*entities.Operation, error) {
	if item, checker := repo.data[id]; checker {
		return item, nil
	}
	return nil, errors.New("элемент не найден")
}

func (repo *RepositoryOperMemory) GetAll() ([]*entities.Operation, error) {
	var result []*entities.Operation
	for _, item := range repo.data {
		result = append(result, item)
	}
	return result, nil
}
