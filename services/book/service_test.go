package book_test

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
	"library/services/book"
	"testing"
)

var bookService book.IBookService
var authorService author.IAuthorService

type bookServiceDsn struct {
	dsn          string
	databaseName string
}

func (dsn bookServiceDsn) GetDsn() string {
	return dsn.dsn
}

func (dsn bookServiceDsn) GetDatabaseName() string {
	return dsn.databaseName
}

func newBookServiceDsn() domain.IDsn {
	databaseName := random.String(10)
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=memory", databaseName)
	return &bookServiceDsn{dsn, databaseName}
}

type bookServiceDb struct {
	db *gorm.DB
}

func (d bookServiceDb) GetDB() *gorm.DB {
	return d.db
}

func newBookServiceDb(dsn domain.IDsn) domain.IDatabase {
	db, err := gorm.Open(sqlite.Open(dsn.GetDsn()), &gorm.Config{})
	if err != nil {
		log.Fatalf("repository [test]: failed to create database: %v\n", err)
	}

	return bookServiceDb{db}
}

func initializeTests() error {
	var err error
	if err = ioc.AddSingleton[domain.IDsn](newBookServiceDsn); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add dsn to service collection")
	}
	if err = ioc.AddTransient[domain.IDatabase](newBookServiceDb); err != nil {
		return errors.Wrap(err, "repository [test]: failed to add database to service collection")
	}
	if err = migrations.NewMigrationRegister().Register(); err != nil {
		return errors.Wrap(err, "repository [test]: failed to register migrations")
	}
	if err = author.NewAuthorServiceRegister().Register(); err != nil {
		return errors.Wrap(err, "repository [test]: failed to register author service")
	}
	if err = book.NewBookServiceRegister().Register(); err != nil {
		return errors.Wrapf(err, "repository [test]: failed to register book service")
	}
	if authorService, err = ioc.Get[author.IAuthorService](); err != nil {
		return errors.Wrap(err, "repository [test]: failed to get author service")
	}
	if bookService, err = ioc.Get[book.IBookService](); err != nil {
		return errors.Wrap(err, "repository [test]: failed to get book service")
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
	migration, err := ioc.Get[migrations.Migration]()
	ass.NoError(err)

	err = migration.CreateDatabase()
	ass.NoError(err)
	err = migration.MigrateDatabase()
	ass.NoError(err)

	return afterTest
}

func afterTest(t *testing.T) {
	ass := assert.New(t)
	migration, err := ioc.Get[migrations.Migration]()
	ass.NoError(err)
	db, err := ioc.Get[domain.IDatabase]()
	ass.NoError(err)
	err = migration.DropTables()
	ass.NoError(err)
	sqlDB, err := db.GetDB().DB()
	ass.NoError(err)
	err = sqlDB.Close()
	ass.NoError(err)
}

func TestBookService_Create(t *testing.T) {
	t.Run("Create book succeeded", createBookSucceeded)
	t.Run("Create book failed when author doesn't exist", createBookFailedWhenAuthorDoesNotExist)
	t.Run("Create book failed when author is empty", createBookFailedWhenAuthorIsEmpty)
	t.Run("Create book failed when author is nil", createBookFailedWhenAuthorIsNil)
	t.Run("Create book failed when title is empty", createBookFailedWhenTitleIsEmpty)
	t.Run("Create book failed when ISBN is too long", createBookFailedWhenISBNIsTooLong)
	t.Run("Create book failed when ISBN is too short", createBookFailedWhenISBNIsTooShort)
	t.Run("Create book failed when trying add same ISBN twice", createBookFailedWhenTryingAddSameISBNTwice)
}

func createBookSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)
	ass.NotNil(b)
}

func createBookFailedWhenAuthorDoesNotExist(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	newBook, err := bookService.Create(uuid.New(),
		random.String(20),
		random.String(13),
		random.String(140),
		uuid.New())

	ass.Error(err)
	ass.Nil(newBook)
}

func createBookFailedWhenISBNIsTooLong(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor(random.String(20))
	ass.Error(err)
	ass.Nil(b)
}

func createBookFailedWhenISBNIsTooShort(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor(random.String(10))
	ass.Error(err)
	ass.Nil(b)
}

func createBookFailedWhenTryingAddSameISBNTwice(t *testing.T) {
	after := beforeTest(t)
	defer after(t)
	ass := assert.New(t)

	isbn := random.String(13)
	b0, err := createBookWithAuthor(isbn)
	ass.NoError(err)
	ass.NotNil(b0)
	ass.Equal(isbn, b0.ISBN)

	b1, err := createBookWithAuthor(isbn)
	ass.Error(err)
	ass.ErrorIs(err, book.ErrBookAlreadyExists)
	ass.Nil(b1)
}

func createBookFailedWhenAuthorIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := bookService.Create(uuid.New(),
		random.String(20),
		random.String(13),
		random.String(140),
		domain.EmptyUUID())
	ass.Error(err)
	ass.Nil(b)
}

func createBookFailedWhenAuthorIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := bookService.Create(uuid.New(),
		random.String(20),
		random.String(13),
		random.String(140),
		uuid.Nil)
	ass.Error(err)
	ass.Nil(b)
}

func createBookFailedWhenTitleIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	bookAuthor, err := createAuthor()
	ass.NoError(err)

	b, err := bookService.Create(uuid.New(),
		"",
		random.String(13),
		random.String(140),
		bookAuthor.ID)
	ass.Error(err)
	ass.Nil(b)
}

func TestBookService_Edit(t *testing.T) {
	t.Run("Update book succeeded", updateBookSucceeded)
	t.Run("Update book failed when author doesn't exist", updateBookFailedWhenAuthorDoesNotExist)
	t.Run("Update book failed when author is empty", updateBookFailedWhenAuthorIsEmpty)
	t.Run("Update book failed when author is nil", updateBookFailedWhenAuthorIsNil)
	t.Run("Update book failed when ISBN is too long", updateBookFailedWhenISBNIsTooLong)
	t.Run("Update book failed when ISBN is too short", updateBookFailedWhenISBNIsTooShort)
}

func updateBookSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	title, description := b.Title, b.Description

	b.Title = random.String(25)
	b.Description = random.String(305)

	b, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.NoError(err)
	ass.NotEqual(title, b.Title)
	ass.NotEqual(description, b.Description)
}

func updateBookFailedWhenAuthorDoesNotExist(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.AuthorID = uuid.New()
	b.Author.ID = b.AuthorID

	b, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
	ass.Nil(b)
}

func updateBookFailedWhenAuthorIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.AuthorID = domain.EmptyUUID()

	b, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
	ass.Nil(b)
}

func updateBookFailedWhenAuthorIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.AuthorID = uuid.Nil

	b, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
	ass.Nil(b)
}

func updateBookFailedWhenISBNIsTooLong(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.ISBN = random.String(20)

	_, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
}

func updateBookFailedWhenISBNIsTooShort(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.ISBN = random.String(10)

	_, err = bookService.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
}

func TestBookService_GetByID(t *testing.T) {
	t.Run("Get book by ID succeeded", getBookByIDSucceeded)
	t.Run("Get book by ID failed when ID is empty", getBookByIDFailedWhenIDIsEmpty)
	t.Run("Get book by ID failed when ID is nil", getBookByIDFailedWhenIDIsNil)
	t.Run("Get book by ID failed when book doesn't exist", getBookByIDFailedWhenBookDoesNotExist)
}

func getBookByIDSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b0, err := createBookWithAuthor("")
	ass.NoError(err)

	b1, err := bookService.GetByID(b0.ID)
	ass.NoError(err)
	ass.Equal(b0.ID, b1.ID)
	ass.Equal(b0.Title, b1.Title)
	ass.Equal(b0.ISBN, b1.ISBN)
	ass.Equal(b0.Description, b1.Description)
	ass.Equal(b0.AuthorID, b1.AuthorID)
}

func getBookByIDFailedWhenIDIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := bookService.GetByID(domain.EmptyUUID())
	ass.Error(err)
}

func getBookByIDFailedWhenIDIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := bookService.GetByID(uuid.Nil)
	ass.Error(err)
}

func getBookByIDFailedWhenBookDoesNotExist(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	_, err := bookService.GetByID(uuid.New())
	ass.Error(err)
}

func TestBookService_GetAll(t *testing.T) {
	t.Run("Get all books succeeded", getAllBooksSucceeded)
	t.Run("Get all books succeeded when there are no book", getAllBooksSucceededWhenThereAreNoBook)
}

func getAllBooksSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	var values []domain.Book
	for i := 0; i < 10; i++ {
		b, err := createBookWithAuthor("")
		ass.NoError(err)
		values = append(values, *b)
	}

	books, err := bookService.GetAll()
	ass.NoError(err)

	ass.Len(books, len(values))

	for _, b := range books {
		r, ok := lo.Find(values, func(val domain.Book) bool {
			return b.ID == val.ID
		})
		ass.True(ok)
		ass.Equal(r.Title, b.Title)
		ass.Equal(r.ISBN, b.ISBN)
		ass.Equal(r.Description, b.Description)
		ass.Equal(r.AuthorID, b.AuthorID)
	}
}

func getAllBooksSucceededWhenThereAreNoBook(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	books, err := bookService.GetAll()
	ass.NoError(err)

	ass.Len(books, 0)
}

func TestBookService_Delete(t *testing.T) {
	t.Run("Delete book succeeded", deleteBookSucceeded)
	t.Run("Delete book failed when ID is empty", deleteBookFailedWhenIDIsEmpty)
	t.Run("Delete book failed when ID is nil", deleteBookFailedWhenIDIsNil)
	t.Run("Delete book failed when book doesn't exist", deleteBookFailedWhenBookDoesNotExist)
}

func deleteBookSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	err = bookService.Delete(b.ID)
	ass.NoError(err)
}

func deleteBookFailedWhenIDIsEmpty(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := bookService.Delete(domain.EmptyUUID())
	ass.Error(err)
}

func deleteBookFailedWhenIDIsNil(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := bookService.Delete(uuid.Nil)
	ass.Error(err)
}

func deleteBookFailedWhenBookDoesNotExist(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	err := bookService.Delete(uuid.New())
	ass.Error(err)
}

func createAuthor() (*domain.Author, error) {
	return authorService.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}

func createBookWithAuthor(isbn string) (*domain.Book, error) {
	bookAuthor, err := addAuthor()
	if err != nil {
		return nil, err
	}

	if isbn == "" {
		return createBook(random.String(13), bookAuthor.ID)
	}

	return createBook(isbn, bookAuthor.ID)
}

func createBook(isbn string, authorID uuid.UUID) (*domain.Book, error) {
	return bookService.Create(uuid.New(),
		random.String(100),
		isbn,
		random.String(240),
		authorID)
}

func addAuthor() (*domain.Author, error) {
	return authorService.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}
