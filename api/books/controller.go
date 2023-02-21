package books

import (
	"github.com/gin-gonic/gin"
	"github.com/mehdihadeli/go-mediatr"
	"library/application/books/commands"
	"library/application/books/queries"
	"library/domain"
	"net/http"
)

type Controller struct {
}

func NewBooksController() *Controller {
	return &Controller{}
}

type createBookDto struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	ISBN        domain.ISBN `json:"isbn"`
	Description string      `json:"description"`
	AuthorID    string      `json:"author_id"`
}

func (a *Controller) Create(ctx *gin.Context) {
	var json createBookDto
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := commands.NewCreateBookCommand(domain.BookID(json.ID),
		json.Title,
		json.ISBN,
		json.Description,
		domain.AuthorID(json.AuthorID))
	response, err := mediatr.Send[*commands.CreateBookCommand, *commands.CreateBookCommandResponse](
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

	query := queries.NewGetBookByIDQuery(domain.BookID(paramID))
	result, err := mediatr.Send[*queries.GetBookByIDQuery, *queries.GetBookByIDQueryResponse](
		ctx.Request.Context(),
		query)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
