package application

import (
	"fmt"
	"library/api/author"
	booksAPI "library/api/books"
	"library/application/authors"
	"library/application/books"
	"library/domain"
	"library/ioc"
	"library/migrations"
	"library/register"
	"library/services"
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
	gormDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(fmt.Printf("App: failed to open database: %v", err))
	}

	db := &database{db: gormDB}
	a.db = db

	if err := ioc.AddSingleton[domain.IDatabase](db); err != nil {
		log.Fatalln(fmt.Printf("App: failed to add database object to IoC: %v", err))
	}

	if err := migrations.CreateAndUseDatabase(dsn); err != nil {
		log.Fatalln(fmt.Printf("App: failed to create and use database: %v", err))
	}

	if err := migrations.UpdateDatabase(); err != nil {
		log.Fatalln(fmt.Printf("App: failed to update database: %v", err))
	}
}

type appRegister struct {
}

func (r *appRegister) Register() error {
	srvRegister := services.NewServiceRegister()
	if err := srvRegister.Register(); err != nil {
		return err
	}

	authorRegister := authors.NewAuthorRegister()
	if err := authorRegister.Register(); err != nil {
		return err
	}

	bookRegister := books.NewBookRegister()
	if err := bookRegister.Register(); err != nil {
		return err
	}

	return nil
}

func (a *App) initializeMediatr() {
	reg := new(appRegister)
	if err := ioc.AddSingleton[register.IRegister[*App]](reg); err != nil {
		log.Fatalln(err)
	}

	if err := reg.Register(); err != nil {
		log.Fatalln(err)
	}
}

func (a *App) initializeRouter() {
	a.router = gin.Default()

	authorCtrl := author.NewAuthorController()
	bookCtrl := booksAPI.NewBooksController()

	v1 := a.router.Group("/api/v1")
	{
		v1.GET("/author", authorCtrl.GetAll)
		v1.GET("/author/:id", authorCtrl.Get)
		v1.POST("/author", authorCtrl.Add)
		v1.PATCH("/author/:id", authorCtrl.Edit)
		v1.DELETE("/author/:id", authorCtrl.Delete)

		v1.POST("/book", bookCtrl.Create)
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
