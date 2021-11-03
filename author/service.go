package author

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"library/models"
)

const (
	FirstNameCol  string = "firstName"
	MiddleNameCol string = "middleName"
	LastNameCol   string = "lastName"
	BooksCol      string = "books"
)

func Add(db *gorm.DB, model *models.Author) error {
	if model.FirstName == "" {
		return errors.New("author: FirstName cannot be empty")
	}
	if model.LastName == "" {
		return errors.New("author: LastName cannot be empty")
	}
	existed, _ := Get(db, model.FirstName, model.LastName)
	if existed != nil {
		return errors.New("author: that author exists")
	}
	result := db.Create(&model)

	return result.Error
}

func Get(db *gorm.DB, firstName, lastName string) (*models.Author, error) {
	if firstName == "" || lastName == "" {
		return nil, errors.New("author: firstName and lastName must be set")
	}

	var author *models.Author
	result := db.Where(fmt.Sprintf("%s = ? AND %s = ?", FirstNameCol, LastNameCol)).Find(&author)

	if result.Error != nil {
		return nil, fmt.Errorf("author: cannot get data: %w", result.Error)
	}

	return author, nil
}
