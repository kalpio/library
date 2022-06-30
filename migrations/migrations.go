package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"library/models"
	"os"
	"strings"

	"gorm.io/gorm"
)

var ErrAutoMigration = errors.New("could not auto migrate")

func CreateAndUseDatabase(db *gorm.DB, name string) error {
	if !strings.Contains(db.Dialector.Name(), "sqlite") {
		if tx := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", name)); tx.Error != nil {
			return tx.Error
		}
		if tx := db.Exec(fmt.Sprintf("USE %s;", name)); tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func DropDatabase(db *gorm.DB, name string) error {
	if !strings.Contains(db.Dialector.Name(), "sqlite") {
		query := fmt.Sprintf(`
USE master;
ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
DROP DATABASE [%s];
`, name, name)

		if tx := db.Exec(query); tx.Error != nil {
			return tx.Error
		}
	} else {
		var (
			database *sql.DB
			err      error
		)
		if database, err = db.DB(); err != nil {
			return err
		}

		if err = database.Close(); err != nil {
			return err
		}

		return os.Remove(fmt.Sprintf("%s.db", name))
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
