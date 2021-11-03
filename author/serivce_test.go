package author_test

import (
	"github.com/matryer/is"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"library/author"
	"library/models"
	"testing"
)

var db *gorm.DB

func init() {
	var err error
	dsn := "sqlserver://jz:jzsoft@localhost:1433?database=library"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Author{}, &models.Book{}); err != nil {
		panic(err)
	}
}

func Test_AddNewAuthor(t *testing.T) {
	a := &models.Author{
		FirstName:  "Robert C.",
		MiddleName: "",
		LastName:   "Martin",
		Books:      nil,
	}

	iss := is.New(t)

	err := author.Add(db, a)

	iss.NoErr(err)
	iss.True(a.ID > 0)
	iss.True(a.FirstName != "")
	iss.True(a.LastName != "")
}
