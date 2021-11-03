package migrations

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"library/models"
)

var ErrAutoMigration = errors.New("could not auto migrate")

func UpdateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Author{}); err != nil {
		return fmt.Errorf("models: %w", ErrAutoMigration)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return fmt.Errorf("book: %w", ErrAutoMigration)
	}

	return nil
}
