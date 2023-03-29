package ioc_test

import (
	"github.com/stretchr/testify/assert"
	"library/ioc"
	"testing"
)

func TestAddSingleton(t *testing.T) {
	t.Parallel()

	t.Run("Add singleton not return error",
		addSingletonNotReturnError)
	t.Run("Add singleton no return error and has registered service",
		addSingletonNotReturnErrorAndHasRegisteredService)
	t.Run("Add singleton succeeded when register same interface twice",
		addSingletonSucceededWhenRegisterSameInterfaceTwice)
	t.Run("Add singleton fail when trying add type which not implement interface",
		addSingletonFailWhenTryingAddTypeWhichNotImplementsInterface)
	t.Run("Add singleton fail when trying add type which is different",
		addSingletonFailWhenTryingAddTypeWhichIsDifferent)
}

func addSingletonNotReturnError(t *testing.T) {
	ass := assert.New(t)
	err := ioc.AddSingleton[iFirstInterface](newFirstImpl)

	ass.NoError(err)
}

func addSingletonNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	err := ioc.AddSingleton[iFirstInterface](newFirstImpl)

	ass.NoError(err)
}

func addSingletonSucceededWhenRegisterSameInterfaceTwice(t *testing.T) {
	ass := assert.New(t)

	err := ioc.AddSingleton[iFirstInterface](newFirstImpl)
	ass.NoError(err)

	err = ioc.AddSingleton[iFirstInterface](newFirstImpl)
	ass.NoError(err)
}

func addSingletonFailWhenTryingAddTypeWhichNotImplementsInterface(t *testing.T) {
	ass := assert.New(t)

	err := ioc.AddSingleton[iFirstInterface](newSecondImpl)
	ass.ErrorIs(err, ioc.ErrInvalidType)
}

func addSingletonFailWhenTryingAddTypeWhichIsDifferent(t *testing.T) {
	ass := assert.New(t)

	err := ioc.AddSingleton[iFirstImpl](newSecondImpl)
	ass.ErrorIs(err, ioc.ErrInvalidType)
}
