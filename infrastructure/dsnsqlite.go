package infrastructure

import (
	"library/domain"
)

func NewDsnSqlite() domain.IDsn {
	return &DsnSqlite{dsn: "library.db"}
}

type DsnSqlite struct {
	dsn string
}

func (d DsnSqlite) GetDsn() string {
	return d.dsn
}
