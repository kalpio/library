package author

import (
	"errors"
	"fmt"
	"library/models"

	"gorm.io/gorm"
)

func Save(db *gorm.DB, author *models.Author) (*models.Author, error) {
	if err := author.Validate(); err != nil {
		return nil, err
	}

	existed, _ := getAuthorByFirstAndLastName(db, author.FirstName, author.LastName)
	if existed != nil && existed.ID > 0 {
		return nil, errors.New("repository: author already exists")
	}

	result := db.Create(&author)

	return author, result.Error
}

func GetByID(db *gorm.DB, id uint) (*models.Author, error) {
	var result *models.Author
	if tx := db.First(&result, id); tx.Error != nil {
		return nil, fmt.Errorf("could not find author by ID: %d: %w", id, tx.Error)
	}

	return result, nil
}

func GetAll(db *gorm.DB) ([]*models.Author, error) {
	var results []*models.Author
	if tx := db.Find(&results); tx.Error != nil {
		return nil, fmt.Errorf("could not read authors: %w", tx.Error)
	}

	return results, nil
}

func getAuthorByFirstAndLastName(db *gorm.DB, firstName, lastName string) (*models.Author, error) {
	if len(firstName) == 0 && len(lastName) == 0 {
		return nil, errors.New("repository: author firstName or lastName must be set")
	}

	var author *models.Author
	result := db.
		Where("firstName = ? AND lastName = ?", firstName, lastName).
		Find(&author)

	if result.Error != nil {
		return nil, fmt.Errorf("repository: cannot get data: %w", result.Error)
	}

	return author, nil
}
