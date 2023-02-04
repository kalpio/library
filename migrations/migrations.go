package migrations

import (
	"database/sql"
	"fmt"
	"library/domain"
	"library/ioc"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var ErrAutoMigration = errors.New("could not auto migrate")

func CreateAndUseDatabase(name string) error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "migrations: failed to get database service")
	}

	if !strings.Contains(db.GetDB().Dialector.Name(), "sqlite") {
		if tx := db.GetDB().Exec(fmt.Sprintf("CREATE DATABASE %s;", name)); tx.Error != nil {
			return tx.Error
		}
		if tx := db.GetDB().Exec(fmt.Sprintf("USE %s;", name)); tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func DropDatabase(name string) error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "migrations: failed to get database service")
	}

	if strings.Contains(getDialectorName(db), "sqlite") {
		return dropSqliteDb(db, name)
	}

	if strings.Contains(getDialectorName(db), "sqlserver") {
		return dropSqlServerDb(db, name)
	}

	return nil
}

func getDialectorName(db domain.IDatabase) string {
	return strings.ToLower(db.GetDB().Dialector.Name())
}

func dropSqliteDb(db domain.IDatabase, name string) error {
	var (
		database *sql.DB
		err      error
	)
	if database, err = db.GetDB().DB(); err != nil {
		return err
	}

	if err = database.Close(); err != nil {
		return err
	}

	return os.Remove(fmt.Sprintf("%s.db", name))
}

func dropSqlServerDb(db domain.IDatabase, name string) error {
	query := fmt.Sprintf(`
USE master;
ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
DROP DATABASE [%s];
`, name, name)

	if tx := db.GetDB().Exec(query); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func UpdateDatabase() error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "migrations: failed to get database service")
	}

	if err := db.GetDB().AutoMigrate(&domain.Author{}); err != nil {
		return fmt.Errorf("models: %w: %v", ErrAutoMigration, err)
	}
	if err := db.GetDB().AutoMigrate(&domain.Book{}); err != nil {
		return fmt.Errorf("book: %w: %v", ErrAutoMigration, err)
	}

	return nil
}
