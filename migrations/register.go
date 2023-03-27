package migrations

import (
	"github.com/pkg/errors"
	"library/ioc"
	"library/register"
)

type migrationRegister struct {
}

func NewMigrationRegister() register.IRegister[*Migration] {
	return &migrationRegister{}
}

func (r *migrationRegister) Register() error {
	err := ioc.AddTransient[Migration](NewMigration)
	return errors.Wrap(err, "register [migration]: failed to register migration")
}
