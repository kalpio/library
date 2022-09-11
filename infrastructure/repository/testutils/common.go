package testutils

import (
	"fmt"
	"library/domain"
	"library/migrations"
	"library/random"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

func (d *database) GetDB() *gorm.DB {
	return d.db
}

func BeforeTest(t *testing.T) (domain.Database, func(t *testing.T)) {
	randomDBName := getRandomDBName()
	var (
		err error
		db  *gorm.DB
	)
	ass := assert.New(t)

	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=rwc", randomDBName)
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	ass.NoError(err)

	ddb := &database{db}

	if err := migrations.CreateAndUseDatabase(ddb, randomDBName); err != nil {
		ass.NoError(err)
	}

	if err := migrations.UpdateDatabase(ddb); err != nil {
		ass.NoError(err)
	}

	return ddb, func(t *testing.T) {
		if err := migrations.DropDatabase(ddb, randomDBName); err != nil {
			ass.NoError(err)
		}
	}
}

func getRandomDBName() string {
	return fmt.Sprintf("library_%s", random.RandomString(6))
}
