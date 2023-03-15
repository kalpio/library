package domain

type IDsn interface {
	GetDsn() string
	GetDatabaseName() string
}
