package testutils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

type repositoryDsn struct {
	dsn          string
	databaseName string
}

func (d repositoryDsn) GetDsn() string {
	return d.dsn
}

func (d repositoryDsn) GetDatabaseName() string {
	return d.databaseName
}

func newRepositoryDsn() domain.IDsn {
	databaseName := getRandomDBName()
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=memory", databaseName)
	return &repositoryDsn{dsn, databaseName}
}

type database struct {
	db *gorm.DB
}

func (d database) GetDB() *gorm.DB {
	return d.db
}

func newRepositoryDatabase(dsn domain.IDsn) domain.IDatabase {
	db, err := gorm.Open(sqlite.Open(dsn.GetDsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("repository [test]: failed to create database: %v\n", err)
	}

	return database{db}
}

func init() {
	if err := ioc.AddSingleton[domain.IDsn](newRepositoryDsn); err != nil {
		log.Fatal("repository [test]: failed to add database DSN to service collection: %v\n", err)
	}

	if err := ioc.AddTransient[domain.IDatabase](newRepositoryDatabase); err != nil {
		log.Fatalf("repository [test]: failed to add database to service collection: %v\n", err)
	}
}

func BeforeTest(t *testing.T) func(t *testing.T) {
	ass := assert.New(t)
	dsn, err := ioc.Get[domain.IDsn]()
	ass.NoError(err)

	err = migrations.CreateAndUseDatabase(dsn.GetDatabaseName())
	ass.NoError(err)
	err = migrations.UpdateDatabase()
	ass.NoError(err)

	return func(t *testing.T) {
		db, err := ioc.Get[domain.IDatabase]()
		ass.NoError(err)
		err = migrations.DropTables()
		ass.NoError(err)
		sqlDB, err := db.GetDB().DB()
		ass.NoError(err)
		err = sqlDB.Close()
		ass.NoError(err)
	}
}

func getRandomDBName() string {
	return fmt.Sprintf("library_%s", random.String(6))
}
