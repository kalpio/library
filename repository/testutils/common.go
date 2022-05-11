package testutils

import (
	"fmt"
	"github.com/matryer/is"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"library/migrations"
	"library/random"
	"testing"
)

func BeforeTest(t *testing.T) (*gorm.DB, func(t *testing.T)) {
	randomDBName := getRandomDBName()
	var (
		err error
		db  *gorm.DB
	)
	iss := is.New(t)

	dsn := "sqlserver://jz:jzsoft@localhost:1433"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	iss.NoErr(err)

	if err := migrations.CreateAndUseDatabase(db, randomDBName); err != nil {
		iss.NoErr(err)
	}

	dsn = fmt.Sprintf("sqlserver://jz:jzsoft@localhost:1433?database=%s", randomDBName)
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	iss.NoErr(err)

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
