package repository

import (
	"github.com/google/uuid"
	"library/domain"
	"library/ioc"

	"github.com/pkg/errors"
)

type Models interface {
	domain.Author | domain.Book
	Validate() error
	GetID() uuid.UUID
}

var (
	errFailedGetDbService = "repository: failed to get database service"
)

func Save[T Models](model T) (T, error) {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return *new(T), errors.Wrap(err, errFailedGetDbService)
	}

	if err := model.Validate(); err != nil {
		return *new(T), errors.Wrapf(err, "repository: an error during %T validation", model)
	}

	tx := db.GetDB().Create(&model)
	return model, errors.Wrapf(tx.Error, "repository: cannot save %T model", model)
}

func Update[T Models](model T) error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, errFailedGetDbService)
	}

	if err := model.Validate(); err != nil {
		return errors.Wrapf(err, "repository: an error during %T validation", model)
	}

	if model.GetID() == domain.EmptyUUID() {
		return errors.Errorf("repository: ID for %T must be set", model)
	}

	if tx := db.GetDB().Model(model).
		Where("id = ?", model.GetID()).
		Updates(model); tx.Error != nil {
		return errors.Wrapf(tx.Error, "repository: cannot update %T model", model)
	}

	return nil
}

func UpdatesInsteadOf[T Models](model T, columns ...string) error {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, errFailedGetDbService)
	}

	if err = model.Validate(); err != nil {
		return errors.Wrapf(err, "repository: an error during %T validate", model)
	}

	if model.GetID() == domain.EmptyUUID() {
		return errors.Errorf("repository: ID for %T must be set", model)
	}

	tx := db.GetDB().Model(model).
		Where("id = ?", model.GetID()).
		Select("*").
		Omit(columns...).
		Updates(model)

	return errors.Wrapf(tx.Error, "repository: cannot update %T model", model)
}

func GetByColumns[T Models](columnValue map[string]interface{}) (T, error) {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return *new(T), errors.Wrap(err, errFailedGetDbService)
	}

	if len(columnValue) == 0 {
		return *new(T), errors.New("repository: column and value cannot be empty")
	}

	var result T
	if tx := db.GetDB().Where(columnValue).Find(&result); tx.Error != nil {
		return *new(T), errors.Wrap(tx.Error, "repository: cannot get data")
	}

	return result, nil
}

func GetAll[T Models]() ([]T, error) {
	var (
		db      domain.IDatabase
		err     error
		results []T
	)

	if db, err = getDB(); err != nil {
		return *new([]T), err
	}

	if err = db.GetDB().Find(&results).Error; err != nil {
		return *new([]T), errors.Wrapf(err, "repository: could not read list %T", new(T))
	}

	return results, nil
}

func GetByID[T Models](id uuid.UUID) (T, error) {
	var (
		db     domain.IDatabase
		err    error
		result T
	)

	db, err = getDB()
	if err != nil {
		return *new(T), err
	}

	if err = db.GetDB().First(&result, id).Error; err != nil {
		return *new(T), errors.Wrapf(err, "repository: could not find %T by ID: %d", new(T), id)
	}

	return result, nil
}

func Delete[T Models](id uuid.UUID) (int64, error) {
	var (
		db  domain.IDatabase
		err error
	)

	db, err = getDB()
	if err != nil {
		return 0, err
	}

	var gormDB = db.GetDB()
	if gormDB = gormDB.Delete(new(T), id); gormDB.Error != nil {
		return 0, errors.Wrapf(gormDB.Error, "repository: could not delete: %T", new(T))
	}

	return gormDB.RowsAffected, nil
}

func getDB() (domain.IDatabase, error) {
	db, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return nil, errors.Wrap(err, errFailedGetDbService)
	}

	return db, nil
}
