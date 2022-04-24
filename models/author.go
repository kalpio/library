package models

import (
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	FirstName  string    `gorm:"column:firstName;index:uq_first_last,unique"`
	MiddleName string    `gorm:"column:middleName"`
	LastName   string    `gorm:"column:lastName;index:uq_first_last,unique"`
	Books      []Book
}
