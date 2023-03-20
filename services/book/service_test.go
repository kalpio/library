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
	t.Run("Create book failed when not author", createBookFailedWhenNoAuthor)
	t.Run("Create book failed when ISBN is too long", createBookFailedWhenISBNIsTooLong)
	t.Run("Create book failed when trying add same ISBN twice", createBookFailedWhenTryingAddSameISBNTwice)
}

func createBookSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	authorSrv, err := ioc.Get[author.IAuthorService]()
	ass.NoError(err)

	bookAuthor, err := addAuthor(authorSrv)

	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	bookData := map[string]interface{}{
		"id":          uuid.New(),
		"title":       random.String(40),
		"isbn":        random.String(13),
		"description": random.String(140),
		"authorID":    bookAuthor.ID,
	}

	newBook, err := bookSrv.Create(bookData["id"].(uuid.UUID),
		bookData["title"].(string),
		bookData["isbn"].(string),
		bookData["description"].(string),
		bookData["authorID"].(uuid.UUID))

	ass.NoError(err)
	ass.Equal(bookData["id"].(uuid.UUID), newBook.ID)
	ass.Equal(bookData["title"].(string), newBook.Title)
	ass.Equal(bookData["isbn"].(string), newBook.ISBN)
	ass.Equal(bookData["description"].(string), newBook.Description)
	ass.Equal(bookData["authorID"].(uuid.UUID), newBook.AuthorID)
}

func createBookFailedWhenNoAuthor(t *testing.T) {
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

	authorSrv, err := ioc.Get[author.IAuthorService]()
	ass.NoError(err)

	bookAuthor, err := addAuthor(authorSrv)
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	newBook, err := bookSrv.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20),
		bookAuthor.ID)

	ass.Error(err)
	ass.Nil(newBook)
}

func createBookFailedWhenTryingAddSameISBNTwice(t *testing.T) {
	after := beforeTest(t)
	defer after(t)
	ass := assert.New(t)

	authorSrv, err := ioc.Get[author.IAuthorService]()
	ass.NoError(err)

	bookAuthor, err := addAuthor(authorSrv)
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	isbn := random.String(13)
	firstBook, err := bookSrv.Create(uuid.New(),
		random.String(20),
		isbn,
		random.String(100),
		bookAuthor.ID)

	ass.NoError(err)
	ass.NotNil(firstBook)

	secondBook, err := bookSrv.Create(uuid.New(),
		random.String(20),
		isbn,
		random.String(100),
		bookAuthor.ID)

	ass.Error(err)
	ass.ErrorIs(err, book.ErrAlreadyExists)
	ass.Nil(secondBook)
}

func TestBookService_Edit(t *testing.T) {
	t.Run("Update book succeeded", updateBookSucceeded)
}

func updateBookSucceeded(t *testing.T) {
	after := beforeTest(t)
	defer after(t)

	ass := assert.New(t)

	ath, err := createAuthor()
	ass.NoError(err)

	isbn := random.String(13)
	b, err := createBook(isbn, ath.ID)
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
