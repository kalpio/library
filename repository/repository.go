package repository

import (
	"bytes"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"library/models"
)

type Models interface {
	models.Author | models.Book
	Validate() error
}

func Save[T Models](db *gorm.DB, model T) error {
	if err := model.Validate(); err != nil {
		return err
	}

	tx := db.Create(model)
	return tx.Error
}

func GetByColumns[T Models](db *gorm.DB, columns []string) (T, error) {
	if len(columns) == 0 {
		return nil, errors.New("repository: columns name cannot be empty")
	}

	var result T
	if tx := db.Where(concatColumnsToWhere(columns)).Find(&result); tx.Error != nil {
		return nil, fmt.Errorf("repository: cannot get data: %w", tx.Error)
	}

	return result, nil
}

func concatColumnsToWhere(columns []string) string {
	length := len(columns)
	var result bytes.Buffer
	for i, column := range columns {
		result.WriteString(fmt.Sprintf("`%s` = ?", column))
		if i != length-1 {
			result.WriteString(", ")
		}
	}

	return result.String()
}

func GetAll[T Models](db *gorm.DB) ([]T, error) {
	var results []T
	if tx := db.Find(&results); tx.Error != nil {
		return nil, fmt.Errorf("repository: could not read %q: %w", new(T), tx.Error)
	}

	return results, nil
}

func GetByID[T Models](db *gorm.DB, id uint) (T, error) {
	var result T
	if tx := db.First(&result, id); tx.Error != nil {
		return nil, fmt.Errorf("repository: could not find %q by ID: %d: %w", new(T), id, tx.Error)
	}

	return result, nil
}
