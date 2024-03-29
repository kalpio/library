package author

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"library/application/authors/commands"
	"library/application/authors/queries"
	"library/domain"
	"net/http"
)

type WebController struct {
}

func NewAuthorWebController() *WebController {
	return &WebController{}
}

type createAuthorFormDto struct {
	FirstName  string `form:"first_name"`
	MiddleName string `form:"middle_name"`
	LastName   string `form:"last_name"`
}

func (a *WebController) Add(ctx *gin.Context) {
	var form createAuthorFormDto

	if err := ctx.Bind(&form); err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	command := commands.NewCreateAuthorCommand(domain.NewAuthorID(), form.FirstName, form.MiddleName, form.LastName)
	response, err := mediatr.Send[*commands.CreateAuthorCommand, *commands.CreateAuthorCommandResponse](
		ctx.Request.Context(),
		command)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusCreated, "author.html", gin.H{"data": response})
}

func (a *WebController) Get(ctx *gin.Context) {
	paramID := ctx.Param("id")

	query := queries.NewGetAuthorByIDQuery(domain.AuthorID(paramID))
	result, err := mediatr.Send[*queries.GetAuthorByIDQuery, *queries.GetAuthorByIDQueryResponse](
		ctx.Request.Context(),
		query)

	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "author.html", gin.H{"data": result})
}

func (a *WebController) GetAll(ctx *gin.Context) {
	query := queries.NewGetAllAuthorsQuery()
	response, err := mediatr.Send[*queries.GetAllAuthorsQuery, *queries.GetAllAuthorsQueryResponse](
		ctx.Request.Context(),
		query)

	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "authors.html", gin.H{"data": response.Result})
}
