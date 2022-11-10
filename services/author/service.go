package author

import (
	"errors"
	"github.com/google/uuid"
	"library/domain"
	"library/infrastructure/repository"
)

type IAuthorService interface {
	Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error)
	Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error)
	GetByID(id uuid.UUID) (*domain.Author, error)
	GetAll() ([]domain.Author, error)
	Delete(id uuid.UUID) (bool, error)
}

type authorService struct {
	db domain.IDatabase
}

func newAuthorService(db domain.IDatabase) IAuthorService {
	return &authorService{db}
}

func (a *authorService) Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := domain.NewAuthor(id, firstName, middleName, lastName)

	exists, err := exists(a.db, firstName, middleName, lastName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("author with that names already exists")
	}

	result, err := repository.Save(a.db, *model)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *authorService) Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := &domain.Author{
		Entity: domain.Entity{
			ID: id,
		},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}

	err := repository.Update(a.db, *model)
	if err != nil {
		return nil, err
	}

	result, err := a.GetByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authorService) GetByID(id uuid.UUID) (*domain.Author, error) {
	result, err := repository.GetByID[domain.Author](a.db, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *authorService) GetAll() ([]domain.Author, error) {
	result, err := repository.GetAll[domain.Author](a.db)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authorService) Delete(id uuid.UUID) (bool, error) {
	var (
		rowsAffected int64
		err          error
	)
	if rowsAffected, err = repository.Delete[domain.Author](a.db, id); err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func exists(db domain.IDatabase, firstName, middleName, lastName string) (bool, error) {
	columns := map[string]interface{}{
		"FirstName":  firstName,
		"MiddleName": middleName,
		"LastName":   lastName,
	}

	result, err := repository.GetByColumns[domain.Author](db, columns)
	if err != nil {
		return false, err
	}

	return len(result.FirstName) > 0 ||
		len(result.MiddleName) > 0 ||
		len(result.LastName) > 0, nil
}
