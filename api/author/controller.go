package author

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authorCtrl struct {
	db *gorm.DB
}

func NewAuthorController(db *gorm.DB) *authorCtrl {
	return &authorCtrl{db: db}
}

func (a *authorCtrl) Add(ctx *gin.Context) {}

func (a *authorCtrl) Get(ctx *gin.Context) {}

func (a *authorCtrl) GetAll(ctx *gin.Context) {}

func (a *authorCtrl) Edit(ctx *gin.Context) {}

func (a *authorCtrl) Delete(ctx *gin.Context) {}
