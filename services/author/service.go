package author

import (
	"github.com/pkg/errors"
	"library/domain"
	"library/infrastructure/repository"

	"github.com/google/uuid"
)

type IAuthorService interface {
	Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error)
	Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error)
	GetByID(id uuid.UUID) (*domain.Author, error)
	GetAll() ([]domain.Author, error)
	Delete(id uuid.UUID) error
}

type authorService struct {
	db domain.IDatabase
}

func newAuthorService(db domain.IDatabase) IAuthorService {
	return &authorService{db}
}

func (a authorService) Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	if a.exists(id) {
		return nil, errors.New("author with that id already exists")
	}

	model := domain.NewAuthor(id, firstName, middleName, lastName)

	result, err := repository.Save(*model)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a authorService) Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := &domain.Author{
		Entity: domain.Entity{
			ID: id,
		},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}

	err := repository.UpdatesInsteadOf(*model, "createdAt")
	if err != nil {
		return nil, err
	}

	result, err := a.GetByID(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a authorService) GetByID(id uuid.UUID) (*domain.Author, error) {
	if id == uuid.Nil {
		return nil, errors.New("author service: author ID must be set")
	}

	result, err := repository.GetByID[domain.Author](id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a authorService) GetAll() ([]domain.Author, error) {
	result, err := repository.GetAll[domain.Author]()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a authorService) Delete(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("author service: author ID must be set")
	}

	var (
		rowsAffected int64
		err          error
	)
	if rowsAffected, err = repository.Delete[domain.Author](id); err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("author service: no rows affected during delete")
	}

	return nil
}

func (a authorService) exists(id uuid.UUID) bool {
	result, err := a.GetByID(id)
	if err != nil {
		return false
	}

	return result.ID == id
}
