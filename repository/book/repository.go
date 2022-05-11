package book

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"library/models"
)

func Save(db *gorm.DB, book *models.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}

	existing, _ := GetByISBN(db, book.ISBN)
	if existing != nil && existing.ID > 0 {
		return errors.New("repository: book already exists")
	}

	tx := db.Create(&book)
	return tx.Error
}

func GetByISBN(db *gorm.DB, isbn string) (*models.Book, error) {
	if len(isbn) == 0 {
		return nil, errors.New("repository: ISBN cannot be empty")
	}
	var result *models.Book
	if tx := db.Where("isbn = ?", isbn).Find(&result); tx.Error != nil {
		return nil, fmt.Errorf("repository: cannot get data: %w", tx.Error)
	}

	return result, nil
}

func GetAll(db *gorm.DB) ([]*models.Book, error) {
	var results []*models.Book
	if tx := db.Find(&results); tx.Error != nil {
		return nil, fmt.Errorf("repository: could not read books: %w", tx.Error)
	}

	return results, nil
}

func Delete(db *gorm.DB, id uint) error {
	if tx := db.Delete(&models.Book{}, id); tx.Error != nil {
		return fmt.Errorf("repository: could not delete book: %w", tx.Error)
	}

	return nil
}
