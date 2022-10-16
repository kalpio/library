package application

import (
	"fmt"
	"library/api/author"
	booksAPI "library/api/books"
	"library/application/authors"
	"library/domain"
	"library/ioc"
	"library/migrations"
	"library/register"
	authorServices "library/services/author"
	bookServices "library/services/book"
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

	if err := ioc.AddSingleton[domain.IDatabase](db); err != nil {
		log.Fatalln(err)
	}
}

type appRegister struct {
}

func (r *appRegister) Register() error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return err
	}

	authorSrv := authorServices.NewAuthorService(db)
	if err := ioc.AddSingleton[authorServices.IAuthorService](authorSrv); err != nil {
		return err
	}

	if err := ioc.AddSingleton[bookServices.IBookService](
		bookServices.NewBookService(db, authorSrv)); err != nil {
		return err
	}

	if err := authors.Register(db); err != nil {
		return err
	}

	reg, err := ioc.Get[register.IRegister[*domain.Book]]()
	if err != nil {
		return err
	}
	reg.Register()

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
