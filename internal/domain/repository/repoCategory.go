package repository

import "task2/internal/domain/entities"

type RepositoryCategory interface {
	Save(category *entities.Category) error
	Delete(id string) error
	Find(id string) (*entities.Category, error)
	GetAll() ([]*entities.Category, error)
}
