package book

import (
	"github.com/pkg/errors"
	"library/domain"
	"library/ioc"
	"library/register"
	"library/services/author"
)

type bookServiceRegister struct {
}

func NewBookServiceRegister() register.IRegister[IBookService] {
	return &bookServiceRegister{}
}

func (r *bookServiceRegister) Register() error {
	database, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "register [book service]: failed to get database service")
	}

	authorSrv, err := ioc.Get[author.IAuthorService]()
	if err != nil {
		return errors.Wrap(err, "register [book service]: failed to get author service")
	}

	bookSrv := newBookService(database, authorSrv)
	if err := ioc.AddSingleton[IBookService](bookSrv); err != nil {
		return errors.Wrap(err, "register [book service]: failed to add book service")
	}

	return nil
}
