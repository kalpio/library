package application

import (
	"fmt"
	"library/api/author"
	booksAPI "library/api/books"
	"library/application/authors"
	"library/application/books"
	"library/domain"
	"library/infrastructure"
	"library/ioc"
	"library/migrations"
	"library/register"
	"library/services"
	"log"
	"net"

	"github.com/gin-gonic/gin"
)

type App struct {
	db     domain.IDatabase
	router *gin.Engine
	host   string
	port   string
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

func (a *App) Initialize() {
	a.initializeDB()
	a.initializeRouter()
	a.initializeMediatr()
}

func (a *App) initializeDB() {
	if err := ioc.AddSingleton[domain.IDsn](infrastructure.NewDsnSqlite); err != nil {
		log.Fatalf("app: failed to add database object to IoC: %v\n", err)
	}

	if err := ioc.AddSingleton[domain.IDatabase](infrastructure.NewDatabase); err != nil {
		log.Fatalf("app: failed to add database object to IoC: %v\n", err)
	}

	dsn, err := ioc.Get[domain.IDsn]()
	if err != nil {
		// TODO: (kalpio) add log
	}
	if err := migrations.CreateAndUseDatabase(dsn.GetDsn()); err != nil {
		log.Fatalf("app: failed to create and use database: %v\n", err)
	}

	if err := migrations.UpdateDatabase(); err != nil {
		log.Fatalf("app: failed to update database: %v\n", err)
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
		v1.GET("/book/:id", bookCtrl.Get)
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
