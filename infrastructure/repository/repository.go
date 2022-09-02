package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"library/domain"
)

type Models interface {
	domain.Author | domain.Book
	Validate() error
}

func Save[T Models](db domain.Database, model T) (T, error) {
	if err := model.Validate(); err != nil {
		return *new(T), err
	}

	tx := db.GetDB().Create(&model)
	return model, tx.Error
}

func GetByColumns[T Models](db domain.Database, columnValue map[string]interface{}) (T, error) {
	if len(columnValue) == 0 {
		return *new(T), errors.New("repository: column and value cannot be empty")
	}

	var result T
	if tx := db.GetDB().Where(columnValue).Find(&result); tx.Error != nil {
		return *new(T), fmt.Errorf("repository: cannot get data: %w", tx.Error)
	}

	return result, nil
}

func GetAll[T Models](db domain.Database) ([]T, error) {
	var results []T
	if tx := db.GetDB().Find(&results); tx.Error != nil {
		return *new([]T), fmt.Errorf("repository: could not read %q: %w", new(T), tx.Error)
	}

	return results, nil
}

func GetByID[T Models](db domain.Database, id uint) (T, error) {
	var result T
	if tx := db.GetDB().First(&result, id); tx.Error != nil {
		return *new(T), fmt.Errorf("repository: could not find %q by ID: %d: %w", new(T), id, tx.Error)
	}

	return result, nil
}

func Delete[T Models](db domain.Database, id uint) (int64, error) {
	var tx *gorm.DB
	if tx = db.GetDB().Delete(new(T), id); tx.Error != nil {
		return 0, fmt.Errorf("repository: could not delete: %w", tx.Error)
	}

	return tx.RowsAffected, nil
}
