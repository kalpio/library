package author

import (
	"library/services/author"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type authorCtrl struct {
	db *gorm.DB
}

func NewAuthorController(db *gorm.DB) *authorCtrl {
	return &authorCtrl{db: db}
}

type authorDto struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func (a *authorCtrl) Add(ctx *gin.Context) {
	var json authorDto

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := author.Create(a.db, json.FirstName, json.MiddleName, json.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (a *authorCtrl) Get(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := author.GetByID(a.db, uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *authorCtrl) GetAll(ctx *gin.Context) {
	result, err := author.GetAll(a.db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *authorCtrl) Edit(ctx *gin.Context) {}

func (a *authorCtrl) Delete(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	succeeded := false
	if succeeded, err = author.Delete(a.db, uint(id)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if succeeded {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{})
}
