package application

import (
	"fmt"
	"library/api/author"
	"library/application/authors"
	"library/application/books"
	"library/domain"
	"library/migrations"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	db     domain.IDatabase
	router *gin.Engine
	host   string
	port   string
}

type database struct {
	db *gorm.DB
}

func (d *database) GetDB() *gorm.DB {
	return d.db
}

func (a *App) DB() domain.IDatabase {
	return a.db
}

func (a *App) Router() *gin.Engine {
	return a.router
}

func (a *App) Host(host string) {
	a.host = host
}

func (a *App) Port(port string) {
	a.port = port
}

func (a *App) Initialize(dsn string) {
	a.initializeDB(dsn)
	a.initializeRouter()
	a.initializeMediatr()
}

func (a *App) initializeDB(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	a.db = &database{db}

	if err := migrations.CreateAndUseDatabase(a.db, dsn); err != nil {
		log.Fatalln(err)
	}
	if err := migrations.UpdateDatabase(a.db); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) initializeMediatr() {
	if err := authors.Register(a.db); err != nil {
		log.Fatalln(err)
	}

	if err := books.Register(a.db); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) initializeRouter() {
	a.router = gin.Default()

	authorCtrl := author.NewAuthorController()

	v1 := a.router.Group("/api/v1")
	{
		v1.GET("/author", authorCtrl.GetAll)
		v1.GET("/author/:id", authorCtrl.Get)
		v1.POST("/author", authorCtrl.Add)
		v1.PATCH("/author/:id", authorCtrl.Edit)
		v1.DELETE("/author/:id", authorCtrl.Delete)
	}
}

func (a *App) Run() {
	if err := a.router.Run(net.JoinHostPort(a.host, a.port)); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) BaseUrl() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(a.host, a.port))
}
