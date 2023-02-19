package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AuthorID string

type Author struct {
	Entity
	FirstName  string `gorm:"column:firstName" json:"first_name"`
	MiddleName string `gorm:"column:middleName" json:"middle_name"`
	LastName   string `gorm:"column:lastName" json:"last_name"`
	Books      []Book `json:"books"`
}

func NewAuthor(id uuid.UUID, firstName, middleName, lastName string) *Author {
	return &Author{
		Entity:     Entity{ID: id},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      nil,
	}
}

func (a Author) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ID, validation.Required),
		validation.Field(&a.FirstName, validation.Required),
		validation.Field(&a.LastName, validation.Required),
	)
}

func (a Author) GetID() uuid.UUID {
	return a.ID
}
