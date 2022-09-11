package domain

import "gorm.io/gorm"

type Database interface {
	GetDB() *gorm.DB
}
