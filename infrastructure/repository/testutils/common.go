package testutils

import (
	"fmt"
	"library/domain"
	"library/ioc"
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

func BeforeTest(t *testing.T) func(t *testing.T) {
	randomDBName := getRandomDBName()
	var (
		err    error
		gormDB *gorm.DB
	)
	ass := assert.New(t)

	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=rwc", randomDBName)
	gormDB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	ass.NoError(err)

	db := &database{gormDB}

	ioc.RemoveSingleton[domain.IDatabase]()
	err = ioc.AddSingleton[domain.IDatabase](db)
	ass.NoError(err)

	if err := migrations.CreateAndUseDatabase(randomDBName); err != nil {
		ass.NoError(err)
	}

	if err := migrations.UpdateDatabase(); err != nil {
		ass.NoError(err)
	}

	return func(t *testing.T) {
		if err := migrations.DropDatabase(randomDBName); err != nil {
			ass.NoError(err)
		}
	}
}

func getRandomDBName() string {
	return fmt.Sprintf("library_%s", random.String(6))
}
