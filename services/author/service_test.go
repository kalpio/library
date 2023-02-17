package author_test

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

	return err
}

func TestAuthorService_Create(t *testing.T) {
	t.Run("Create author succeeded when all fields provided", createAuthorSucceededWhenAllFieldsProvided)
	t.Run("Create author fail when first name is empty", createAuthorFailWhenFirstNameIsEmpty)
	t.Run("Create author fail when last name is empty", createAuthorFailWhenLastNameIsEmpty)
	t.Run("Create author succeeded when middle name is empty", createAuthorSucceededWhenMiddleNameIsEmpty)
	t.Run("Create author fail when try to add authors with same ID", createAuthorFailWhenTryToAddAuthorsWithSameID)
	t.Run("Create author succeeded when try to add authors with same data and different ID",
		createAuthorSucceededWhenTryToAddAuthorsWithSameDataAndDifferentID)
}

func TestAuthorService_Delete(t *testing.T) {
	t.Run("Delete author succeeded when ID is set", deleteAuthorSucceededWhenIDIsSet)
	t.Run("Delete author failed when ID isn't set", deleteAuthorFailedWhenIDIsNotSet)
	t.Run("Delete author failed when ID doesn't exists", deleteAuthorFailedWhenIDDoesNotExist)
}

func TestAuthorService_Edit(t *testing.T) {
	t.Run("Edit author succeeded when ID is set", nil)
	t.Run("Edit author fail when ID isn't set", nil)
	t.Run("Edit author fail when author doesn't exist", nil)
	t.Run("Edit author fail when first name is empty", nil)
	t.Run("Edit author fail when last name is empty", nil)
	t.Run("Edit author succeeded when middle name is empty", nil)
}

func createAuthorSucceededWhenAllFieldsProvided(t *testing.T) {
	ass := assert.New(t)

	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  random.String(20),
		"middleName": random.String(20),
		"lastName":   random.String(20),
	}

	newAuthor, err := executeServiceCreateAuthor(data)
	ass.NoError(err)
	assertAuthor(ass, data, newAuthor)
}

func createAuthorFailWhenFirstNameIsEmpty(t *testing.T) {
	ass := assert.New(t)

	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  "",
		"middleName": random.String(20),
		"lastName":   random.String(20),
	}

	_, err := executeServiceCreateAuthor(data)
	ass.Error(err)
}

func createAuthorFailWhenLastNameIsEmpty(t *testing.T) {
	ass := assert.New(t)

	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  random.String(20),
		"middleName": random.String(20),
		"lastName":   "",
	}

	_, err := executeServiceCreateAuthor(data)
	ass.Error(err)
}

func createAuthorSucceededWhenMiddleNameIsEmpty(t *testing.T) {
	ass := assert.New(t)

	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  random.String(20),
		"middleName": "",
		"lastName":   random.String(20),
	}

	ath, err := executeServiceCreateAuthor(data)
	ass.NoError(err)
	assertAuthor(ass, data, ath)
}

func createAuthorFailWhenTryToAddAuthorsWithSameID(t *testing.T) {
	ass := assert.New(t)
	id := uuid.New()

	data := map[string]interface{}{
		"id":         id,
		"firstName":  random.String(20),
		"middleName": random.String(20),
		"lastName":   random.String(20),
	}

	_, err := executeServiceCreateAuthor(data)
	ass.NoError(err)

	data = map[string]interface{}{
		"id":         id,
		"firstName":  random.String(20),
		"middleName": random.String(20),
		"lastName":   random.String(20),
	}

	_, err = executeServiceCreateAuthor(data)
	ass.Error(err)
}

func createAuthorSucceededWhenTryToAddAuthorsWithSameDataAndDifferentID(t *testing.T) {
	ass := assert.New(t)

	firstName := random.String(20)
	middleName := random.String(20)
	lastName := random.String(20)

	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  firstName,
		"middleName": middleName,
		"lastName":   lastName,
	}

	ath0, err := executeServiceCreateAuthor(data)
	ass.NoError(err)
	assertAuthor(ass, data, ath0)

	data = map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  firstName,
		"middleName": middleName,
		"lastName":   lastName,
	}

	ath1, err := executeServiceCreateAuthor(data)
	ass.NoError(err)
	assertAuthor(ass, data, ath1)
}

func createFakeAuthor() (*domain.Author, error) {
	data := map[string]interface{}{
		"id":         uuid.New(),
		"firstName":  random.String(20),
		"middleName": random.String(20),
		"lastName":   random.String(20),
	}

	return executeServiceCreateAuthor(data)
}

func executeServiceCreateAuthor(data map[string]interface{}) (*domain.Author, error) {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, err
	}

	return authorSrv.Create(data["id"].(uuid.UUID),
		data["firstName"].(string),
		data["middleName"].(string),
		data["lastName"].(string))
}

func assertAuthor(ass *assert.Assertions, data map[string]interface{}, ath *domain.Author) {
	ass.Equal(data["id"].(uuid.UUID), ath.ID)
	ass.Equal(data["firstName"].(string), ath.FirstName)
	ass.Equal(data["middleName"].(string), ath.MiddleName)
	ass.Equal(data["lastName"].(string), ath.LastName)
}

func deleteAuthorSucceededWhenIDIsSet(t *testing.T) {
	ass := assert.New(t)
	ath, err := createFakeAuthor()
	ass.NoError(err)

	err = executeServiceDeleteAuthor(ath.ID)
	ass.NoError(err)
}

func deleteAuthorFailedWhenIDIsNotSet(t *testing.T) {
	ass := assert.New(t)
	err := executeServiceDeleteAuthor(uuid.Nil)
	ass.Error(err)
}

func deleteAuthorFailedWhenIDDoesNotExist(t *testing.T) {
	ass := assert.New(t)
	err := executeServiceDeleteAuthor(uuid.New())
	ass.Error(err)
}

func executeServiceDeleteAuthor(id uuid.UUID) error {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return err
	}

	return authorSrv.Delete(id)
}

func executeServiceEditAuthor(data map[string]interface{}) (*domain.Author, error) {
	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return nil, err
	}

	return authorSrv.Edit(data["id"].(uuid.UUID),
		data["firstName"].(string),
		data["middleName"].(string),
		data["lastName"].(string))
}
