package author_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"library/domain"
	"library/ioc"
	"library/migrations"
	"library/random"
	"library/services/author"
	"testing"
)

type authorServiceDsn struct {
	dsn          string
	databaseName string
}

func (dsn authorServiceDsn) GetDsn() string {
	return dsn.dsn
}

func (dsn authorServiceDsn) GetDatabaseName() string {
	return dsn.databaseName
}

func newAuthorServiceDsn() domain.IDsn {
	databaseName := random.String(10)
	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=memory", databaseName)
	return &authorServiceDsn{dsn, databaseName}
}

type authorServiceDb struct {
	db *gorm.DB
}

func (d authorServiceDb) GetDB() *gorm.DB {
	return d.db
}

func newAuthorServiceDb(dsn domain.IDsn) domain.IDatabase {
	db, err := gorm.Open(sqlite.Open(dsn.GetDsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("repository [test]: failed to create database: %v\n", err)
	}
	return authorServiceDb{db}
}

func initializeTests() error {
	if err := ioc.AddSingleton[domain.IDsn](newAuthorServiceDsn); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add dsn to container")
	}
	if err := ioc.AddTransient[domain.IDatabase](newAuthorServiceDb); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add database to container")
	}
	if err := author.NewAuthorServiceRegister().Register(); err != nil {
		return errors.Wrap(err, "repository [test]: failed to register author service")
	}
	return nil
}

func init() {
	if err := initializeTests(); err != nil {
		log.Fatalf("repository [test]: failed to initialize tests: %v\n", err)
	}
}

func beforeTest(t *testing.T) func(t *testing.T) {
	ass := assert.New(t)
	dsn, err := ioc.Get[domain.IDsn]()
	ass.NoError(err)

	err = migrations.CreateAndUseDatabase(dsn.GetDatabaseName())
	ass.NoError(err)
	err = migrations.UpdateDatabase()
	ass.NoError(err)

	return afterTest
}

func afterTest(t *testing.T) {
	ass := assert.New(t)
	db, err := ioc.Get[domain.IDatabase]()
	ass.NoError(err)
	err = migrations.DropTables()
	ass.NoError(err)
	sqlDB, err := db.GetDB().DB()
	ass.NoError(err)
	err = sqlDB.Close()
	ass.NoError(err)
}

func TestAuthorService_Create(t *testing.T) {
	t.Run("create author succeeded", createAuthorSucceeded)
	t.Run("create author failed when already exists", createAuthorFailedWhenAlreadyExists)
	t.Run("create author failed when first name is empty", createAuthorFailedWhenFirstNameIsEmpty)
	t.Run("create author failed when last name is empty", createAuthorFailedWhenLastNameIsEmpty)
	t.Run("create author failed when id is empty", createAuthorFailedWhenIdIsEmpty)
	t.Run("create author failed when id is nil", createAuthorFailedWhenIdIsNil)
}

func createAuthorFailedWhenIdIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithValues(uuid.Nil, random.String(10), random.String(10), random.String(20))
	ass.Error(err)
}

func createAuthorFailedWhenIdIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithValues(domain.EmptyUUID(), random.String(10), random.String(10), random.String(20))
	ass.Error(err)
}

func createAuthorFailedWhenLastNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithNames(random.String(10), random.String(10), "")
	ass.Error(err)
}

func createAuthorFailedWhenFirstNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithNames("", random.String(10), random.String(20))
	ass.Error(err)
}

func createAuthorSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := createAuthor()
	ass.NoError(err)
	ass.NotNil(a)
}

func createAuthorFailedWhenAlreadyExists(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := createAuthor()
	ass.NoError(err)

	_, err = createAuthorWithNames(a.FirstName, a.MiddleName, a.LastName)
	ass.Error(err)
}

func createAuthor() (*domain.Author, error) {
	return createAuthorWithNames(random.String(10), random.String(10), random.String(20))
}

func createAuthorWithNames(firstName, middleName, lastName string) (*domain.Author, error) {
	return createAuthorWithValues(uuid.New(), firstName, middleName, lastName)
}

func createAuthorWithValues(id uuid.UUID, firstName, middleName, lastName string) (*domain.Author, error) {
	authorService, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, errors.Wrap(err, "repository [test]: failed to get author service")
	}
	return authorService.Create(id, firstName, middleName, lastName)
}
