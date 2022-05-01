package migrations

import (
	"errors"
	"fmt"
	"library/models"

	"gorm.io/gorm"
)

var ErrAutoMigration = errors.New("could not auto migrate")

func CreateAndUseDatabase(db *gorm.DB, name string) error {
	if tx := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", name)); tx.Error != nil {
		return tx.Error
	}
	if tx := db.Exec(fmt.Sprintf("USE %s;", name)); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func DropDatabase(db *gorm.DB, name string) error {
	sql := fmt.Sprintf(`
USE master;
ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
DROP DATABASE [%s];
`, name, name)

	if tx := db.Exec(sql); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func UpdateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Author{}); err != nil {
		return fmt.Errorf("models: %w: %v", ErrAutoMigration, err)
	}
	if err := db.AutoMigrate(&models.Book{}); err != nil {
		return fmt.Errorf("book: %w: %v", ErrAutoMigration, err)
	}

	return nil
}
