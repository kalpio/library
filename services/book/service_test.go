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
	t.Run("Create book failed when ISBN is too long", createBookFailedWhenISBNIsTooLong)
	t.Run("Create book failed when trying add the same ISBN twice", createBookFailedWhenTryingAddTheSameISBNTwice)
	t.Run("Edit book succeeded", editBookSucceeded)
	t.Run("Edit book failed when no author", editBookFailedWhenNoAuthor)
	t.Run("Edit book failed when ISBN is too long", editBookFailedWhenISBNIsTooLong)
	t.Run("Edit book failed when title is empty", editBookFailedWhenTitleIsEmpty)
	t.Run("Delete book succeeded", deleteBookSucceeded)
	t.Run("Delete book failed when book ID is empty", deleteBookFailedWhenBookIDIsEmpty)
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

func createBookFailedWhenISBNIsTooLong(t *testing.T) {
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

func createBookFailedWhenTryingAddTheSameISBNTwice(t *testing.T) {
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

func editBookSucceeded(t *testing.T) {
	ass := assert.New(t)

	b, err := createFakeBook()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	editedBook, err := bookSrv.Edit(b.ID,
		b.Title,
		b.ISBN,
		b.Description,
		b.AuthorID)
	ass.NoError(err)
	ass.Equal(b.ID, editedBook.ID)
	ass.Equal(b.Title, editedBook.Title)
	ass.Equal(b.ISBN, editedBook.ISBN)
	ass.Equal(b.Description, editedBook.Description)
	ass.Equal(b.AuthorID, editedBook.AuthorID)
}

func editBookFailedWhenNoAuthor(t *testing.T) {
	ass := assert.New(t)

	b, err := createFakeBook()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	_, err = bookSrv.Edit(b.ID,
		b.Title,
		b.ISBN,
		b.Description,
		uuid.Nil)
	ass.Error(err)
}

func editBookFailedWhenISBNIsTooLong(t *testing.T) {
	ass := assert.New(t)

	b, err := createFakeBook()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	_, err = bookSrv.Edit(b.ID,
		b.Title,
		domain.ISBN(random.String(30)),
		b.Description,
		b.AuthorID)
	ass.Error(err)
}

func editBookFailedWhenTitleIsEmpty(t *testing.T) {
	ass := assert.New(t)

	b, err := createFakeBook()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	_, err = bookSrv.Edit(b.ID,
		"",
		b.ISBN,
		b.Description,
		b.AuthorID)
	ass.Error(err)
}

func deleteBookSucceeded(t *testing.T) {
	ass := assert.New(t)

	b, err := createFakeBook()
	ass.NoError(err)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	err = bookSrv.Delete(b.ID)
	ass.NoError(err)

	_, err = bookSrv.GetByID(b.ID)
	ass.Error(err)
}

func deleteBookFailedWhenBookIDIsEmpty(t *testing.T) {
	ass := assert.New(t)

	bookSrv, err := ioc.Get[book.IBookService]()
	ass.NoError(err)

	err = bookSrv.Delete(uuid.Nil)
	ass.Error(err)
}

func createFakeBook() (*domain.Book, error) {
	ath, err := createFakeAuthor()
	if err != nil {
		return nil, err
	}

	bookSrv, err := ioc.Get[book.IBookService]()
	if err != nil {
		return nil, err
	}

	return bookSrv.Create(uuid.New(),
		random.String(20),
		domain.ISBN(random.String(13)),
		random.String(100),
		ath.ID)
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
