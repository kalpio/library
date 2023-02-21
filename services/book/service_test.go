package book_test

import (
	"library/domain"
	"library/ioc"
	"library/random"
	"library/services/author"
	"library/services/book"
	"testing"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type testDB struct {
	db *gorm.DB
}

func (t testDB) GetDB() *gorm.DB {
	return t.db
}

func init() {
	if err := initializeTests(); err != nil {
		log.Fatal(err)
	}
}

func newDB() (testDB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return testDB{}, err
	}

	if err = db.AutoMigrate(&domain.Author{}, &domain.Book{}); err != nil {
		return testDB{}, err
	}

	return testDB{db: db}, nil
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

func TestBookService_Create(t *testing.T) {
	t.Run("Create book succeeded", createBookSucceeded)
	t.Run("Create book failed when no author", createBookFailedWhenNoAuthor)
	t.Run("Create failed when ISBN is too long", createFailedWhenISBNIsTooLong)
	t.Run("Create failed when trying add the same ISBN twice", createFailedWhenTryingAddTheSameISBNTwice)
}

func createBookSucceeded(t *testing.T) {
	ass := assert.New(t)

	bookAuthor, err := createFakeAuthor()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	bookData := map[string]interface{}{
		"id":          uuid.New(),
		"title":       random.String(40),
		"isbn":        domain.ISBN(random.String(13)),
		"description": random.String(140),
		"authorID":    bookAuthor.ID,
	}

	newBook, err := bookSrv.Create(bookData["id"].(uuid.UUID),
		bookData["title"].(string),
		bookData["isbn"].(domain.ISBN),
		bookData["description"].(string),
		bookData["authorID"].(uuid.UUID))

	ass.NoError(err)
	ass.Equal(bookData["id"].(uuid.UUID), newBook.ID)
	ass.Equal(bookData["title"].(string), newBook.Title)
	ass.Equal(bookData["isbn"].(domain.ISBN), newBook.ISBN)
	ass.Equal(bookData["description"].(string), newBook.Description)
	ass.Equal(bookData["authorID"].(uuid.UUID), newBook.AuthorID)
}

func createBookFailedWhenNoAuthor(t *testing.T) {
	ass := assert.New(t)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	newBook, err := bookSrv.Create(uuid.New(),
		random.String(20),
		domain.ISBN(random.String(13)),
		random.String(140),
		uuid.New())

	ass.Error(err)
	ass.Nil(newBook)
}

func createFailedWhenISBNIsTooLong(t *testing.T) {
	ass := assert.New(t)

	bookAuthor, err := createFakeAuthor()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	newBook, err := bookSrv.Create(uuid.New(),
		random.String(20),
		domain.ISBN(random.String(20)),
		random.String(20),
		bookAuthor.ID)

	ass.Error(err)
	ass.Nil(newBook)
}

func createFailedWhenTryingAddTheSameISBNTwice(t *testing.T) {
	ass := assert.New(t)

	bookAuthor, err := createFakeAuthor()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	isbn := domain.ISBN(random.String(13))
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

func createFakeAuthor() (*domain.Author, error) {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, err
	}

	return authorSrv.Create(uuid.New(),
		random.String(20),
		random.String(20),
		random.String(20))
}
