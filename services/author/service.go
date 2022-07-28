package author

import (
	"errors"
	"gorm.io/gorm"
	"library/models"
	"library/repository"
)

func Create(db *gorm.DB, firstName, middleName, lastName string) (*models.Author, error) {
	model := &models.Author{
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

func exists(db *gorm.DB, firstName, middleName, lastName string) (bool, error) {
	columns := map[string]interface{}{
		"FirstName":  firstName,
		"MiddleName": middleName,
		"LastName":   lastName,
	}

	result, err := repository.GetByColumns[models.Author](db, columns)
	if err != nil {
		return false, err
	}

	return len(result.FirstName) > 0 ||
		len(result.MiddleName) > 0 ||
		len(result.LastName) > 0, nil
}
