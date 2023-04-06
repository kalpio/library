package bookstest

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"library/application/books"
	"library/domain"
	"library/ioc"
	"library/random"
	"library/services/book"
)

type dsnBook struct {
	dsn          string
	databaseName string
}

func (d dsnBook) GetDsn() string {
	return d.dsn
}

func (d dsnBook) GetDatabaseName() string {
	return d.databaseName
}

func newDsnBook() domain.IDsn {
	return dsnBook{dsn: "", databaseName: ""}
}

type dbBook struct {
}

func newDBBook(_ domain.IDsn) domain.IDatabase {
	return dbBook{}
}

func (d dbBook) GetDB() *gorm.DB {
	return nil
}

func (d dbBook) GetDatabaseName() string {
	return ""
}

func init() {
	if err := ioc.AddSingleton[domain.IDsn](newDsnBook); err != nil {
		log.Fatalln(err)
	}

	if err := ioc.AddSingleton[domain.IDatabase](newDBBook); err != nil {
		log.Fatalln(err)
	}

	if err := ioc.AddSingleton[book.IBookService](newBookServiceMock); err != nil {
		log.Fatalln(err)
	}

	bookRegister := books.NewBookRegister()
	if err := bookRegister.Register(); err != nil {
		log.Fatalln(err)
	}
}

type BookServiceMock struct {
	mock.Mock
}

func newBookServiceMock() *BookServiceMock {
	return &BookServiceMock{}
}

func (b *BookServiceMock) Create(id domain.BookID,
	title, isbn, description string,
	authorID domain.AuthorID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) Edit(id domain.BookID,
	title, isbn, description string,
	authorID domain.AuthorID) (*domain.Book, error) {

	args := b.Called(id, title, isbn, description, authorID)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) GetByID(id domain.BookID) (*domain.Book, error) {
	args := b.Called(id)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (b *BookServiceMock) GetAll() ([]domain.Book, error) {
	args := b.Called()
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (b *BookServiceMock) Delete(id domain.BookID) error {
	args := b.Called(id)
	return args.Error(0)
}

func CreateBook() *domain.Book {
	return domain.NewBook(domain.NewBookID(),
		random.String(20),
		random.String(12),
		random.String(120),
		domain.NewAuthor(domain.NewAuthorID(),
			random.String(20),
			random.String(20),
			random.String(20)))
}
