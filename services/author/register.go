package author

import (
	"github.com/pkg/errors"
	"library/ioc"
	"library/register"
)

type authorServiceRegister struct {
}

func NewAuthorServiceRegister() register.IRegister[IAuthorService] {
	return &authorServiceRegister{}
}

func (r *authorServiceRegister) Register() error {
	if err := ioc.AddTransient[IAuthorService](newAuthorService); err != nil {
		return errors.Wrap(err, "register [author service]: failed to add author service")
	}

	return nil
}
