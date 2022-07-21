package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"library/api/author"
	"log"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	_ = author.NewAuthorController(db)
}
