package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	FirstName  string `gorm:"column:firstName;index:uq_first_last,unique"`
	MiddleName string `gorm:"column:middleName"`
	LastName   string `gorm:"column:lastName;index:uq_first_last,unique"`
	Books      []Book
}

func NewAuthor(firstName, middleName, lastName string) *Author {
	return &Author{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      []Book{},
	}
}

func (a Author) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName, validation.Required),
		validation.Field(&a.LastName, validation.Required),
	)
}
