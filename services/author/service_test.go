package author_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
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

var authorService author.IAuthorService

type authorServiceDsn struct {
	dsn          string
	databaseName string
}

func newAuthorServiceDsn() domain.IDsn {
	databaseName := random.String(10)
	dsn := fmt.Sprintf("file:%s.db?cache=shared&mode=memory", databaseName)
	return &authorServiceDsn{dsn, databaseName}
}

func (dsn authorServiceDsn) GetDsn() string {
	return dsn.dsn
}

func (dsn authorServiceDsn) GetDatabaseName() string {
	return dsn.databaseName
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
	var err error
	if err = ioc.AddSingleton[domain.IDsn](newAuthorServiceDsn); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add dsn to container")
	}
	if err = ioc.AddTransient[domain.IDatabase](newAuthorServiceDb); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add database to container")
	}
	if err = author.NewAuthorServiceRegister().Register(); err != nil {
		return errors.Wrap(err, "repository [test]: failed to register author service")
	}
	if authorService, err = ioc.Get[author.IAuthorService](); err != nil {
		return errors.Wrap(err, "repository [test]: failed to get author service")
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

func createAuthorSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := createAuthor()
	ass.NoError(err)
	ass.NotNil(a)
}

func createAuthorFailedWhenFirstNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithNames("", random.String(10), random.String(20))
	ass.Error(err)
}

func createAuthorFailedWhenLastNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithNames(random.String(10), random.String(10), "")
	ass.Error(err)
}

func createAuthorFailedWhenIdIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithValues(domain.EmptyUUID(), random.String(10), random.String(10), random.String(20))
	ass.Error(err)
}

func createAuthorFailedWhenIdIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := createAuthorWithValues(uuid.Nil, random.String(10), random.String(10), random.String(20))
	ass.Error(err)
}

func TestAuthorService_Edit(t *testing.T) {
	t.Run("edit author succeeded", editAuthorSucceeded)
	t.Run("edit author failed when not exists", editAuthorFailedWhenNotExists)
	t.Run("edit author failed when first name is empty", editAuthorFailedWhenFirstNameIsEmpty)
	t.Run("edit author failed when last name is empty", editAuthorFailedWhenLastNameIsEmpty)
	t.Run("edit author failed when id is empty", editAuthorFailedWhenIdIsEmpty)
	t.Run("edit author failed when id is nil", editAuthorFailedWhenIdIsNil)
}

func editAuthorSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a0.FirstName = random.String(10)
	a0.MiddleName = random.String(10)
	a0.LastName = random.String(20)

	a1, err := authorService.Edit(a0.ID, a0.FirstName, a0.MiddleName, a0.LastName)
	ass.NoError(err)

	ass.Equal(a0.ID, a1.ID)
	ass.Equal(a0.FirstName, a1.FirstName)
	ass.Equal(a0.MiddleName, a1.MiddleName)
	ass.Equal(a0.LastName, a1.LastName)
}

func editAuthorFailedWhenNotExists(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.Edit(uuid.New(), a0.FirstName, a0.MiddleName, a0.LastName)
	ass.Error(err)
	ass.Nil(a1)
}

func editAuthorFailedWhenFirstNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.Edit(a0.ID, "", a0.MiddleName, a0.LastName)
	ass.Error(err)
	ass.Nil(a1)
}

func editAuthorFailedWhenLastNameIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.Edit(a0.ID, a0.FirstName, a0.MiddleName, "")
	ass.Error(err)
	ass.Nil(a1)
}

func editAuthorFailedWhenIdIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.Edit(domain.EmptyUUID(), a0.FirstName, a0.MiddleName, a0.LastName)
	ass.Error(err)
	ass.Nil(a1)
}

func editAuthorFailedWhenIdIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.Edit(uuid.Nil, a0.FirstName, a0.MiddleName, a0.LastName)
	ass.Error(err)
	ass.Nil(a1)
}

func TestAuthorService_GetByID(t *testing.T) {
	t.Run("get by id succeeded", getByIDSucceeded)
	t.Run("get by id failed when not exists", getByIDFailedWhenNotExists)
	t.Run("get by id failed when id is empty", getByIDFailedWhenIDIsEmpty)
	t.Run("get by id failed when id is nil", getByIDFailedWhenIDIsNil)
}

func getByIDSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a0, err := createAuthor()
	ass.NoError(err)

	a1, err := authorService.GetByID(a0.ID)
	ass.NoError(err)
	ass.Equal(a0.ID, a1.ID)
	ass.Equal(a0.FirstName, a1.FirstName)
	ass.Equal(a0.MiddleName, a1.MiddleName)
	ass.Equal(a0.LastName, a1.LastName)
}

func getByIDFailedWhenNotExists(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := authorService.GetByID(uuid.New())
	ass.Error(err)
	ass.Nil(a)
}

func getByIDFailedWhenIDIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := authorService.GetByID(domain.EmptyUUID())
	ass.Error(err)
	ass.Nil(a)
}

func getByIDFailedWhenIDIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := authorService.GetByID(uuid.Nil)
	ass.Error(err)
	ass.Nil(a)
}

func TestAuthorService_GetAll(t *testing.T) {
	t.Run("get all succeeded", getAllSucceeded)
	t.Run("get all succeeded when there are no authors", getAllSucceededWhenThereAreNoAuthors)
}

func getAllSucceededWhenThereAreNoAuthors(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	authors, err := authorService.GetAll()
	ass.NoError(err)
	ass.Len(authors, 0)
}

func TestAuthorService_Delete(t *testing.T) {
	t.Run("delete succeeded", deleteSucceeded)
	t.Run("delete failed when not exists", deleteFailedWhenNotExists)
	t.Run("delete failed when id is empty", deleteFailedWhenIDIsEmpty)
	t.Run("delete failed when id is nil", deleteFailedWhenIDIsNil)
}

func deleteSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	a, err := createAuthor()
	ass.NoError(err)

	err = authorService.Delete(a.ID)
	ass.NoError(err)
}

func deleteFailedWhenNotExists(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := authorService.Delete(uuid.New())
	ass.Error(err)
}

func deleteFailedWhenIDIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := authorService.Delete(domain.EmptyUUID())
	ass.Error(err)
}

func deleteFailedWhenIDIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := authorService.Delete(uuid.Nil)
	ass.Error(err)
}

func getAllSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	var values []*domain.Author
	for i := 0; i < 10; i++ {
		a, err := createAuthor()
		ass.NoError(err)
		values = append(values, a)
	}

	authors, err := authorService.GetAll()
	ass.NoError(err)
	ass.Len(authors, len(values))

	for _, v := range values {
		a, ok := lo.Find(authors, func(val domain.Author) bool {
			return val.ID == v.ID
		})
		ass.True(ok)
		ass.Equal(v.FirstName, a.FirstName)
		ass.Equal(v.MiddleName, a.MiddleName)
		ass.Equal(v.LastName, a.LastName)
	}
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
	return authorService.Create(id, firstName, middleName, lastName)
}
