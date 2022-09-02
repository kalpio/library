package author

import (
	"errors"
	"library/domain"
	"library/infrastructure/repository"
)

func Create(db domain.Database, firstName, middleName, lastName string) (*domain.Author, error) {
	model := &domain.Author{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}

	exists, err := exists(db, firstName, middleName, lastName)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("author with that names alredy exists")
	}

	result, err := repository.Save(db, *model)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetByID(db domain.Database, id uint) (*domain.Author, error) {
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

func Delete(db domain.Database, id uint) (bool, error) {
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
