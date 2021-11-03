package models

import (
	"gorm.io/gorm"
	"time"
)

type Author struct {
	gorm.Model
	FirstName  string    `gorm:"column:firstName;index:uq_first_last,unique"`
	MiddleName string    `gorm:"column:middleName"`
	LastName   string    `gorm:"column:lastName;index:uq_first_last,unique"`
	BirthDate  time.Time `gorm:"column:birthDate"`
	Books      []Book
}
