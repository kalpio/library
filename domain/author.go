package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type AuthorID string

func (a AuthorID) String() string {
	return string(a)
}

func (a AuthorID) IsNil() bool {
	return a == AuthorID(uuid.Nil.String())
}

func (a AuthorID) IsEmpty() bool {
	return a == AuthorID(EmptyUUID().String())
}

func (a AuthorID) UUID() uuid.UUID {
	return uuid.MustParse(a.String())
}

func NewAuthorID() AuthorID {
	return AuthorID(uuid.NewString())
}

type Author struct {
	Entity[AuthorID]
	FirstName  string `gorm:"column:firstName;index:uq_first_last,unique" json:"first_name"`
	MiddleName string `gorm:"column:middleName" json:"middle_name"`
	LastName   string `gorm:"column:lastName;index:uq_first_last,unique" json:"last_name"`
	Books      []Book `json:"books"`
}

func NewAuthor(id AuthorID, firstName, middleName, lastName string) *Author {
	return &Author{
		Entity:     Entity[AuthorID]{ID: id},
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
	return a.ID.UUID()
}

func (a Author) validateID(_ interface{}) error {
	if a.ID.IsNil() {
		return validation.NewError("id", "id is null")
	}
	if a.ID.IsEmpty() {
		return validation.NewError("id", "id is empty")
	}

	return nil
}
