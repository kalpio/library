package book

import (
	"github.com/pkg/errors"
	"library/ioc"
	"library/register"
)

type bookServiceRegister struct {
}

func NewBookServiceRegister() register.IRegister[IBookService] {
	return &bookServiceRegister{}
}

func (r *bookServiceRegister) Register() error {
	if err := ioc.AddSingleton[IBookService](newBookService); err != nil {
		return errors.Wrap(err, "register [book service]: failed to add book service")
	}

	return nil
}
