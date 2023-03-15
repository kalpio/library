package application

import (
	"fmt"
	"library/domain"
	"library/infrastructure"
	"library/ioc"
	"library/migrations"
	"library/register"
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
	a.configureServices()
	a.migrateDatabase()
	a.router = configureRouter()
	a.initializeMediatr()
}

func (*App) configureServices() {
	if err := ioc.AddSingleton[domain.IDsn](infrastructure.NewDsnSqlite); err != nil {
		log.Fatalf("app: failed to add database DSN to service collection: %v\n", err)
	}

	if err := ioc.AddSingleton[domain.IDatabase](infrastructure.NewDatabase); err != nil {
		log.Fatalf("app: failed to add database to service collection: %v\n", err)
	}

	if err := ioc.AddSingleton[register.IRegister[*App]](newRegister); err != nil {
		log.Fatalf("app: failed to add register to service collection: %v\n", err)
	}
}

func (*App) migrateDatabase() {
	dsn, err := ioc.Get[domain.IDsn]()
	if err != nil {
		log.Fatalf("app: failed to get DNS service instance: %v\n", err)
	}

	if err := migrations.CreateAndUseDatabase(dsn.GetDatabaseName()); err != nil {
		log.Fatalf("app: failed to create and use database: %v\n", err)
	}

	if err := migrations.UpdateDatabase(); err != nil {
		log.Fatalf("app: failed to update database: %v\n", err)
	}
}

func (a *App) initializeMediatr() {
	reg, err := ioc.Get[register.IRegister[*App]]()
	if err != nil {
		log.Fatalf("app: failed to get register service instance: %v\n", err)
	}

	if err := reg.Register(); err != nil {
		log.Fatalf("app: failed to register mediatr services: %v\n", err)
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
