package application

import (
	"github.com/pkg/errors"
	"library/application/authors"
	"library/application/books"
	"library/migrations"
	"library/register"
	"library/services"
)

type reg struct {
}

func newRegister() register.IRegister[*App] {
	return new(reg)
}

func (reg) Register() error {
	migrationRegister := migrations.NewMigrationRegister()
	if err := migrationRegister.Register(); err != nil {
		return errors.Wrap(err, "register [app]: failed to register migration")
	}

	serviceRegister := services.NewServiceRegister()
	if err := serviceRegister.Register(); err != nil {
		return errors.Wrap(err, "register [app]: failed to register services")
	}

	authorsRegister := authors.NewAuthorRegister()
	if err := authorsRegister.Register(); err != nil {
		return errors.Wrap(err, "register [app]: failed to register authors")
	}

	booksRegister := books.NewBookRegister()
	if err := booksRegister.Register(); err != nil {
		return errors.Wrap(err, "register [app]: failed to register books")
	}

	return nil
}
