package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"library/domain"
)

type database struct {
	db  *gorm.DB
	dsn domain.IDsn
}

func NewDatabase(dsn domain.IDsn) domain.IDatabase {
	db, err := gorm.Open(sqlite.Open(dsn.GetDsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("app: failed to open database: %v", err)
	}

	return database{db: db, dsn: dsn}
}

func (d database) GetDB() *gorm.DB {
	return d.db
}

func (d database) GetDatabaseName() string {
	return d.dsn.GetDatabaseName()
}
