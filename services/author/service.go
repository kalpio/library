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

func (a *authorService) Create(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := domain.NewAuthor(domain.AuthorID(id.String()), firstName, middleName, lastName)

	exists, err := exists(firstName, middleName, lastName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("author with that names already exists")
	}

	result, err := repository.Save(*model)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *authorService) Edit(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := &domain.Author{
		Entity: domain.Entity[domain.AuthorID]{
			ID: domain.AuthorID(id.String()),
		},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}

	err := repository.Update(*model)
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
	result, err := repository.GetByID[domain.Author](id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *authorService) GetAll() ([]domain.Author, error) {
	result, err := repository.GetAll[domain.Author]()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *authorService) Delete(id uuid.UUID) error {
	err := repository.Delete[domain.Author](id)
	return errors.Wrap(err, "could not delete author")
}

func exists(firstName, middleName, lastName string) (bool, error) {
	columns := map[string]interface{}{
		"FirstName":  firstName,
		"MiddleName": middleName,
		"LastName":   lastName,
	}

	result, err := repository.GetByColumns[domain.Author](columns)
	if err != nil {
		return false, err
	}

	return len(result.FirstName) > 0 ||
		len(result.MiddleName) > 0 ||
		len(result.LastName) > 0, nil
}
