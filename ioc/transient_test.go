package ioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTransient(t *testing.T) {
	t.Run("Add transient not return error",
		addTransientNotReturnError)
	t.Run("Add transient not return error and has registered service",
		addTransientNotReturnErrorAndHasRegisteredService)
}

func addTransientNotReturnError(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddTransient[iFirstInterface](newFirstImpl)
	ass.NoError(err)
}

func addTransientNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddTransient[iFirstInterface](newFirstImpl)

	ass.NoError(err)
	ass.Len(values, 1)
}
