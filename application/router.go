package application

import (
	"library/api/author"
	"library/api/books"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func configureRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "hx-current-url", "hx-request", "hx-target", "hx-trigger"},
		MaxAge:       12 * time.Hour,
	}))

	r.Static("/assets", "./templates/assets")
	r.LoadHTMLFiles(
		"templates/index.html",
		"templates/error.html",
		"templates/author.html",
		"templates/authors.html")

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

	authorWebCtrl := author.NewAuthorWebController()

	web := r.Group("/")
	{
		web.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", gin.H{})
		})
		web.GET("/author", authorWebCtrl.GetAll)
		web.GET("/author/:id", authorWebCtrl.Get)
		web.POST("/author", authorWebCtrl.Add)
	}

	return r
}
