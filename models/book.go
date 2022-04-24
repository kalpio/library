package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string
	ISBN     string
	Content  []byte
	Format   string
	Version  string
	AuthorID uint
	Author   Author
}
