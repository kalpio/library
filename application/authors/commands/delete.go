package commands

type DeleteAuthorCommand struct {
	AuthorID uint
}

func NewDeleteAuthorCommand(authorID uint) *DeleteAuthorCommand {
	return &DeleteAuthorCommand{AuthorID: authorID}
}

type DeleteAuthorCommandResponse struct {
	Succeeded bool
}
