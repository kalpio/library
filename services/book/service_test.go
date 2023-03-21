package book_test

import (
	"fmt"
	"github.com/google/uuid"
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
	if err := ioc.AddSingleton[domain.IDsn](newBookServiceDsn); err != nil {
		log.Fatalf("repository [test]: failed to add database DSN to service collection: %v\n", err)
	}

	if err := ioc.AddTransient[domain.IDatabase](newBookServiceDb); err != nil {
		log.Fatalf("repository [test]: failed to add database to service collection: %v\n", err)
	}

	if err := author.NewAuthorServiceRegister().Register(); err != nil {
		return err
	}

	if err := book.NewBookServiceRegister().Register(); err != nil {
		return err
	}

	return nil
}

func init() {
	if err := initializeTests(); err != nil {
		log.Fatal(err)
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	newBook, err := bookSrv.Create(uuid.New(),
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err := bookSrv.Create(uuid.New(),
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err := bookSrv.Create(uuid.New(),
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err := bookSrv.Create(uuid.New(),
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	b, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
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

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	_, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
}

func updateBookFailedWhenISBNIsTooShort(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	b, err := createBookWithAuthor("")
	ass.NoError(err)

	b.ISBN = random.String(10)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	_, err = bookSrv.Edit(b.ID, b.Title, b.ISBN, b.Description, b.AuthorID)
	ass.Error(err)
}

func createAuthor() (*domain.Author, error) {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, err
	}

	return authorSrv.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}

func createBookWithAuthor(isbn string) (*domain.Book, error) {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, err
	}

	bookAuthor, err := addAuthor(authorSrv)
	if err != nil {
		return nil, err
	}

	if isbn == "" {
		return createBook(random.String(13), bookAuthor.ID)
	}

	return createBook(isbn, bookAuthor.ID)
}

func createBook(isbn string, authorID uuid.UUID) (*domain.Book, error) {
	bookSrv, err := ioc.Get[book.IBookService]()
	if err != nil {
		return nil, err
	}

	return bookSrv.Create(uuid.New(),
		random.String(100),
		isbn,
		random.String(240),
		authorID)
}

func addAuthor(authorSrv author.IAuthorService) (*domain.Author, error) {
	return authorSrv.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}
