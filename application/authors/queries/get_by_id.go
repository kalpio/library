package queries

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"library/domain"
	"time"
)

type GetAuthorByIDQuery struct {
	AuthorID domain.AuthorID
}

type GetAuthorByIDQueryResponse struct {
	AuthorID   uuid.UUID      `json:"id"`
	FirstName  string         `json:"first_name"`
	MiddleName string         `json:"middle_name"`
	LastName   string         `json:"last_name"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func NewGetAuthorByIDQuery(authorID domain.AuthorID) *GetAuthorByIDQuery {
	return &GetAuthorByIDQuery{AuthorID: authorID}
}
