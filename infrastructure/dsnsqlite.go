package infrastructure

import (
	"fmt"
	"library/domain"
)

func NewDsnSqlite() domain.IDsn {
	var databaseName = "library"
	var dsn = fmt.Sprintf("file:%s.db?cache=shared&more=rwc", databaseName)
	return &DsnSqlite{dsn, databaseName}
}

type DsnSqlite struct {
	dsn          string
	databaseName string
}

func (d DsnSqlite) GetDsn() string {
	return d.dsn
}

func (d DsnSqlite) GetDatabaseName() string {
	return d.databaseName
}
