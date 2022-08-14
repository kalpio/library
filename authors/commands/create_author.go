package commands

type CreateAuthorCommand struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CreateAuthorCommandResponse struct {
	AuthorID uint
}

func NewCreateAuthorCommand(firstName, middleName, lastName string) *CreateAuthorCommand {
	return &CreateAuthorCommand{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
