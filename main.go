package main

import (
	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	authorCtrl := author.NewAuthorController(db)

	v1 := router.Group("/v1")
	{
		v1.GET("/author", authorCtrl.GetAll)
		v1.GET("/author/:id", authorCtrl.Get)
		v1.POST("/author", authorCtrl.Add)
		v1.PATCH("/author/:id", authorCtrl.Edit)
		v1.DELETE("/author/:id", authorCtrl.Delete)
	}

	router.Run(":8089")
}
