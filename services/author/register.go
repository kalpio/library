package author

import (
	"github.com/pkg/errors"
	"library/domain"
	"library/ioc"
	"library/register"
)

type authorServiceRegister struct {
}

func NewAuthorServiceRegister() register.IRegister[IAuthorService] {
	return &authorServiceRegister{}
}

func (r *authorServiceRegister) Register() error {
	database, err := ioc.Get[domain.IDatabase]()
	if err != nil {
		return errors.Wrap(err, "register [author service]: failed to get database service")
	}

	authorSrv := newAuthorService(database)
	if err := ioc.AddSingleton[IAuthorService](authorSrv); err != nil {
		return errors.Wrap(err, "register [author service]: failed to add author service")
	}

	return nil
}
