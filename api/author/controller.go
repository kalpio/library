package author

import (
	"library/application/authors/commands"
	"library/application/authors/queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
)

type authorCtrl struct {
}

func NewAuthorController() *authorCtrl {
	return &authorCtrl{}
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

	command := commands.NewCreateAuthorCommand(json.FirstName, json.MiddleName, json.LastName)
	response, err := mediatr.Send[*commands.CreateAuthorCommand, *commands.CreateAuthorCommandResponse](ctx.Request.Context(), command)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (a *authorCtrl) Get(ctx *gin.Context) {
	paramID := ctx.Param("id")
	authorID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := queries.NewGetAuthorByIDQuery(uint(authorID))
	result, err := mediatr.Send[*queries.GetAuthorByIDQuery, *queries.GetAuthorByIDQueryResponse](ctx.Request.Context(), query)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *authorCtrl) GetAll(ctx *gin.Context) {
	query := queries.NewGetAllAuthorsQuery()
	response, err := mediatr.Send[*queries.GetAllAuthorsQuery, *queries.GetAllAuthorsQueryResponse](ctx.Request.Context(), query)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Result)
}

func (a *authorCtrl) Edit(_ *gin.Context) {}

func (a *authorCtrl) Delete(ctx *gin.Context) {
	paramID := ctx.Param("id")
	authorID, err := strconv.ParseUint(paramID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := commands.NewDeleteAuthorCommand(uint(authorID))
	response, err := mediatr.Send[*commands.DeleteAuthorCommand, *commands.DeleteAuthorCommandResponse](ctx.Request.Context(), command)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if response.Succeeded {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	ctx.JSON(http.StatusBadRequest, gin.H{})
}
