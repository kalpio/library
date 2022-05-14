package repository

import (
	"fmt"
	"gorm.io/gorm"
	"library/models"
)

type Models interface {
	models.Author | models.Book
}

func GetAll[T Models](db *gorm.DB) ([]*T, error) {
	var results []*T
	if tx := db.Find(&results); tx.Error != nil {
		return nil, fmt.Errorf("repository: could not read %s: %w", new(T), tx.Error)
	}

	return results, nil
}
