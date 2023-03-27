package migrations

import (
	"database/sql"
	"fmt"
	"library/domain"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const sqlServerDialectorName = "sqlserver"
const sqliteDialectorName = "sqlite"

type Migration struct {
	db domain.IDatabase
}

func NewMigration(db domain.IDatabase) Migration {
	return Migration{db: db}
}

func (m Migration) CreateDatabase() error {
	databaseName := m.getDatabaseName()
	if !strings.Contains(m.getDialectorName(), "sqlite") {
		if tx := m.db.GetDB().Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName)); tx.Error != nil {
			return tx.Error
		}
		if tx := m.db.GetDB().Exec(fmt.Sprintf("USE %s;", databaseName)); tx.Error != nil {
			return tx.Error
		}
	}

	return nil
}

func (m Migration) DropDatabase() error {
	if strings.Contains(m.getDialectorName(), sqliteDialectorName) {
		return m.dropSqliteDatabase()
	}

	if strings.Contains(m.getDialectorName(), sqlServerDialectorName) {
		return m.dropSqlServerDatabase()
	}

	return nil
}

func (m Migration) MigrateDatabase() error {
	if err := m.db.GetDB().AutoMigrate(&domain.Author{}); err != nil {
		return errors.Wrap(err, "migrations: failed to migrate author table")
	}
	if err := m.db.GetDB().AutoMigrate(&domain.Book{}); err != nil {
		return errors.Wrap(err, "migrations: failed to migrate book table")
	}

	return nil
}

func (m Migration) DropTables() error {
	if err := m.db.GetDB().Migrator().DropTable(&domain.Author{}); err != nil {
		return errors.Wrap(err, "migrations: failed to drop author table")
	}
	if err := m.db.GetDB().Migrator().DropTable(&domain.Book{}); err != nil {
		return errors.Wrap(err, "migrations: failed to drop book table")
	}

	return nil
}

func (m Migration) getDatabaseName() string {
	return m.db.GetDatabaseName()
}

func (m Migration) getDialectorName() string {
	return strings.ToLower(m.db.GetDB().Dialector.Name())
}

func (m Migration) dropSqliteDatabase() error {
	databaseName := m.getDatabaseName()

	var (
		database *sql.DB
		err      error
	)
	if database, err = m.db.GetDB().DB(); err != nil {
		return err
	}
	if err = database.Close(); err != nil {
		return err
	}

	return os.Remove(fmt.Sprintf("%s.db", databaseName))
}

func (m Migration) dropSqlServerDatabase() error {
	databaseName := m.getDatabaseName()
	query := fmt.Sprintf(`
USE master;
ALTER DATABASE [%s] SET SINGLE_USER WITH ROLLBACK IMMEDIATE;
DROP DATABASE [%s];
`, databaseName, databaseName)

	if tx := m.db.GetDB().Exec(query); tx.Error != nil {
		return tx.Error
	}

	return nil
}
