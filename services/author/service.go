package author

import (
	"errors"
	"github.com/google/uuid"
	"library/domain"
	"library/infrastructure/repository"
)

func Create(db domain.Database, id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := domain.NewAuthor(id, firstName, middleName, lastName)

	exists, err := exists(db, firstName, middleName, lastName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("author with that names already exists")
	}

	result, err := repository.Save(db, *model)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func Edit(db domain.Database, id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	model := &domain.Author{
		Entity: domain.Entity{
			ID: id,
		},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}

	err := repository.Update(db, *model)
	if err != nil {
		return nil, err
	}

	result, err := GetByID(db, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetByID(db domain.Database, id uuid.UUID) (*domain.Author, error) {
	result, err := repository.GetByID[domain.Author](db, id)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetAll(db domain.Database) ([]domain.Author, error) {
	result, err := repository.GetAll[domain.Author](db)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Delete(db domain.Database, id uuid.UUID) (bool, error) {
	var (
		rowsAffected int64
		err          error
	)
	if rowsAffected, err = repository.Delete[domain.Author](db, id); err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func exists(db domain.Database, firstName, middleName, lastName string) (bool, error) {
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
