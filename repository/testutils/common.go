package testutils

import (
	"fmt"
	"library/migrations"
	"library/random"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func BeforeTest(t *testing.T) (*gorm.DB, func(t *testing.T)) {
	randomDBName := getRandomDBName()
	var (
		err error
		db  *gorm.DB
	)
	ass := assert.New(t)

	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=rwc", randomDBName)
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	ass.NoError(err)

	if err := migrations.CreateAndUseDatabase(db, randomDBName); err != nil {
		ass.NoError(err)
	}

	if err := migrations.UpdateDatabase(db); err != nil {
		ass.NoError(err)
	}

	return db, func(t *testing.T) {
		if err := migrations.DropDatabase(db, randomDBName); err != nil {
			ass.NoError(err)
		}
	}
}

func getRandomDBName() string {
	return fmt.Sprintf("library_%s", random.RandomString(6))
}
