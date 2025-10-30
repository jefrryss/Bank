package repository

import "task2/internal/domain/entities"

type RepositoryOperations interface {
	Save(operation *entities.Operation) error
	Delete(id string) error
	Find(id string) (*entities.Operation, error)
	GetAll() ([]*entities.Operation, error)
}
