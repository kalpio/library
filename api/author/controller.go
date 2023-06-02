package author

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"library/application/authors/commands"
	"library/application/authors/queries"
	"library/domain"
	"net/http"
)

type Controller struct {
}

func NewAuthorController() *Controller {
	return &Controller{}
}

type createAuthorDto struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func (a *Controller) Add(ctx *gin.Context) {
	var json createAuthorDto

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := commands.NewCreateAuthorCommand(domain.NewAuthorID(), json.FirstName, json.MiddleName, json.LastName)
	response, err := mediatr.Send[*commands.CreateAuthorCommand, *commands.CreateAuthorCommandResponse](
		ctx.Request.Context(),
		command)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (a *Controller) Get(ctx *gin.Context) {
	paramID := ctx.Param("id")

	query := queries.NewGetAuthorByIDQuery(domain.AuthorID(paramID))
	result, err := mediatr.Send[*queries.GetAuthorByIDQuery, *queries.GetAuthorByIDQueryResponse](
		ctx.Request.Context(),
		query)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (a *Controller) GetAll(ctx *gin.Context) {
	query := queries.NewGetAllAuthorsQuery()
	response, err := mediatr.Send[*queries.GetAllAuthorsQuery, *queries.GetAllAuthorsQueryResponse](
		ctx.Request.Context(),
		query)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.Result)
}

type editAuthorDto struct {
	ID string `json:"id"`
	createAuthorDto
}

func (a *Controller) Edit(ctx *gin.Context) {
	var json editAuthorDto

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := commands.NewEditAuthorCommand(domain.AuthorID(json.ID), json.FirstName, json.MiddleName, json.LastName)
	response, err := mediatr.Send[*commands.EditAuthorCommand, *commands.EditAuthorCommandResponse](
		ctx.Request.Context(),
		command)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *Controller) Delete(ctx *gin.Context) {
	paramID := ctx.Param("id")

	command := commands.NewDeleteAuthorCommand(domain.AuthorID(paramID))
	response, err := mediatr.Send[*commands.DeleteAuthorCommand, *commands.DeleteAuthorCommandResponse](
		ctx.Request.Context(),
		command)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !response.Succeeded {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Unable to delete author"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
