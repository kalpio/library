package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Author struct {
	Entity
	FirstName  string `gorm:"column:firstName;index:uq_first_last,unique" json:"first_name"`
	MiddleName string `gorm:"column:middleName" json:"middle_name"`
	LastName   string `gorm:"column:lastName;index:uq_first_last,unique" json:"last_name"`
	Books      []Book `json:"books"`
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

func (a Author) GetID() uint {
	return a.ID
}
