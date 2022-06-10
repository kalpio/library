package repository

import (
	"errors"
	"fmt"
	"library/models"

	"gorm.io/gorm"
)

type Models interface {
	models.Author | models.Book
	Validate() error
}

func Save[T Models](db *gorm.DB, model T) (T, error) {
	if err := model.Validate(); err != nil {
		return *new(T), err
	}

	tx := db.Create(&model)
	return model, tx.Error
}

func GetByColumns[T Models](db *gorm.DB, columnValue map[string]interface{}) (T, error) {
	if len(columnValue) == 0 {
		return *new(T), errors.New("repository: column and value cannot be empty")
	}

	var result T
	if tx := db.Where(columnValue).Find(&result); tx.Error != nil {
		return *new(T), fmt.Errorf("repository: cannot get data: %w", tx.Error)
	}

	return result, nil
}

func GetAll[T Models](db *gorm.DB) ([]T, error) {
	var results []T
	if tx := db.Find(&results); tx.Error != nil {
		return *new([]T), fmt.Errorf("repository: could not read %q: %w", new(T), tx.Error)
	}

	return results, nil
}

func GetByID[T Models](db *gorm.DB, id uint) (T, error) {
	var result T
	if tx := db.First(&result, id); tx.Error != nil {
		return *new(T), fmt.Errorf("repository: could not find %q by ID: %d: %w", new(T), id, tx.Error)
	}

	return result, nil
}

func Delete[T Models](db *gorm.DB, id uint) error {
	if tx := db.Delete(new(T), id); tx.Error != nil {
		return fmt.Errorf("repository: could not delete: %w", tx.Error)
	}

	return nil
}
