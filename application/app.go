package application

import (
	"fmt"
	"library/api/author"
	"library/migrations"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
	host   string
	port   string
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
}

func (a *App) initializeDB(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	a.DB = db

	if err := migrations.CreateAndUseDatabase(a.DB, dsn); err != nil {
		log.Fatalln(err)
	}
	if err := migrations.UpdateDatabase(a.DB); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) initializeRouter() {
	a.Router = gin.Default()

	authorCtrl := author.NewAuthorController(a.DB)

	v1 := a.Router.Group("/api/v1")
	{
		v1.GET("/author", authorCtrl.GetAll)
		v1.GET("/author/:id", authorCtrl.Get)
		v1.POST("/author", authorCtrl.Add)
		v1.PATCH("/author/:id", authorCtrl.Edit)
		v1.DELETE("/author/:id", authorCtrl.Delete)
	}
}

func (a *App) Run() {
	if err := a.Router.Run(net.JoinHostPort(a.host, a.port)); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) BaseUrl() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(a.host, a.port))
}