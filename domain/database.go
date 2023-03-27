package domain

import "gorm.io/gorm"

type IDatabase interface {
	GetDB() *gorm.DB
	GetDatabaseName() string
}
