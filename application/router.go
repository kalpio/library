package application

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"library/api/author"
	"library/api/books"
	"time"
)

func configureRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "hx-current-url", "hx-request", "hx-target", "hx-trigger"},
		MaxAge:       12 * time.Hour,
	}))

	authorCtrl := author.NewAuthorController()
	bookCtrl := books.NewBooksController()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/author", authorCtrl.GetAll)
		v1.GET("/author/:id", authorCtrl.Get)
		v1.POST("/author", authorCtrl.Add)
		v1.PATCH("/author/:id", authorCtrl.Edit)
		v1.DELETE("/author/:id", authorCtrl.Delete)

		v1.GET("/book", bookCtrl.GetAll)
		v1.GET("/book/:id", bookCtrl.Get)
		v1.POST("/book", bookCtrl.Add)
		v1.PATCH("/book/:id", bookCtrl.Edit)
		v1.DELETE("/book/:id", bookCtrl.Delete)
	}

	return r
}
