package services

import (
	"library/register"
	"library/services/author"
	"library/services/book"
)

type IService interface {
}

type servicesRegister struct {
}

func NewServiceRegister() register.IRegister[IService] {
	return &servicesRegister{}
}

func (r *servicesRegister) Register() error {
	var lastErr error

	authorSrvRegister := author.NewAuthorServiceRegister()
	if err := authorSrvRegister.Register(); err != nil {
		lastErr = err
	}

	bookSrvRegister := book.NewBookServiceRegister()
	if err := bookSrvRegister.Register(); err != nil {
		lastErr = err
	}

	return lastErr
}
