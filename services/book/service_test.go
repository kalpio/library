package book_test

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"library/domain"
	"library/ioc"
	"library/random"
	"library/services/author"
	"library/services/book"
	"testing"
)

type testDB struct {
	db *gorm.DB
}

func (t *testDB) GetDB() *gorm.DB {
	return t.db
}

func newDB() (*testDB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&domain.Author{}); err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&domain.Book{}); err != nil {
		return nil, err
	}

	return &testDB{db: db}, nil
}

func initializeTests() error {
	db, err := newDB()
	if err != nil {
		return err
	}
	if err := ioc.AddSingleton[domain.IDatabase](db); err != nil {
		return err
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

func TestBookService_CreateBookSucceeded(t *testing.T) {
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

func TestBookService_CreateFail_WhenNoAuthor(t *testing.T) {
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

func TestBookService_CreateFail_WhenISBN_IsTooLong(t *testing.T) {
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

func TestBookService_CreateFail_WhenTryingAddSameISBNTwice(t *testing.T) {
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

func addAuthor(authorSrv author.IAuthorService) (*domain.Author, error) {
	return authorSrv.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}
