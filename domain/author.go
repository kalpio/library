package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AuthorID string

func ParseUUID[T AuthorID | BookID](val T) uuid.UUID {
	return uuid.MustParse(string(val))
}

func ParseID[T AuthorID | BookID](val uuid.UUID) T {
	return T(val.String())
}

type Author struct {
	Entity
	FirstName  string `gorm:"column:firstName;index:uq_first_last,unique" json:"first_name"`
	MiddleName string `gorm:"column:middleName" json:"middle_name"`
	LastName   string `gorm:"column:lastName;index:uq_first_last,unique" json:"last_name"`
	Books      []Book `json:"books"`
}

func NewAuthor(id uuid.UUID, firstName, middleName, lastName string) *Author {
	return &Author{
		Entity:     Entity{ID: id},
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Books:      []Book{},
	}
}

func (a Author) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ID, validation.By(a.validateID)),
		validation.Field(&a.FirstName, validation.Required),
		validation.Field(&a.LastName, validation.Required),
	)
}

func (a Author) GetID() uuid.UUID {
	return a.ID
}

func (a Author) validateID(_ interface{}) error {
	if a.ID == uuid.Nil {
		return validation.NewError("id", "id is null")
	}
	if a.ID == EmptyUUID() {
		return validation.NewError("id", "id is empty")
	}

	return nil
}
