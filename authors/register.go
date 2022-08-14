package authors

import (
	"library/authors/commands"

	mediatr "github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) {
	createAuthorCommandHandler := commands.NewCreateAuthorCommandHandler(db)
	mediatr.RegisterRequestHandler[*commands.CreateAuthorCommand, *commands.CreateAuthorCommandResponse](createAuthorCommandHandler)
}
