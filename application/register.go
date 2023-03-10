package application

import (
	"library/application/authors"
	"library/application/books"
	"library/register"
	"library/services"
)

type reg struct {
}

func newRegister() register.IRegister[*App] {
	return new(reg)
}

func (reg) Register() error {
	serviceRegister := services.NewServiceRegister()
	if err := serviceRegister.Register(); err != nil {
		return err
	}

	authorsRegister := authors.NewAuthorRegister()
	if err := authorsRegister.Register(); err != nil {
		return err
	}

	booksRegister := books.NewBookRegister()
	if err := booksRegister.Register(); err != nil {
		return err
	}

	return nil
}
