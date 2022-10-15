package application

import (
	"library/domain"
)

type IRegister interface {
	Register(db domain.IDatabase) error
}
