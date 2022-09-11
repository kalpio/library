package commands

type CreateAuthorCommand struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type CreateAuthorCommandResponse struct {
	AuthorID   uint   `json:"id"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

func NewCreateAuthorCommand(firstName, middleName, lastName string) *CreateAuthorCommand {
	return &CreateAuthorCommand{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}
}
