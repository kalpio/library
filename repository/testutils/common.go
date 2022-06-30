package testutils

import (
	"fmt"
	"library/migrations"
	"library/random"
	"testing"

	"github.com/matryer/is"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func BeforeTest(t *testing.T) (*gorm.DB, func(t *testing.T)) {
	randomDBName := getRandomDBName()
	var (
		err error
		db  *gorm.DB
	)
	iss := is.New(t)

	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=rwc", randomDBName)
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	iss.NoErr(err)

	if err := migrations.CreateAndUseDatabase(db, randomDBName); err != nil {
		iss.NoErr(err)
	}

	if err := migrations.UpdateDatabase(db); err != nil {
		iss.NoErr(err)
	}

	return db, func(t *testing.T) {
		if err := migrations.DropDatabase(db, randomDBName); err != nil {
			iss.NoErr(err)
		}
	}
}

func getRandomDBName() string {
	return fmt.Sprintf("library_%s", random.RandomString(6))
}
