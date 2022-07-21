package author

import "gorm.io/gorm"

type authorCtrl struct {
	db *gorm.DB
}

func NewAuthorController(db *gorm.DB) *authorCtrl {
	return &authorCtrl{db: db}
}
