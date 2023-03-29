package ioc_test

import (
	"library/ioc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTransient(t *testing.T) {
	t.Parallel()

	t.Run("Add transient not return error",
		addTransientNotReturnError)
	t.Run("Add transient not return error and has registered service",
		addTransientNotReturnErrorAndHasRegisteredService)
}

func addTransientNotReturnError(t *testing.T) {
	ass := assert.New(t)
	err := ioc.AddTransient[iFirstInterface](newFirstImpl)
	ass.NoError(err)
}

func addTransientNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	err := ioc.AddTransient[iFirstInterface](newFirstImpl)

	ass.NoError(err)
}
