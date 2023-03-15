package ioc

import (
	"github.com/stretchr/testify/assert"
	"library/random"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("Get return registered singleton instance",
		getReturnRegisteredSingletonInstance)
	t.Run("Get return not unique registered singleton instance",
		getReturnNotUniqueRegisteredSingletonInstance)
	t.Run("Get return registered singleton pointer",
		getReturnRegisteredSingletonPointer)
	t.Run("Get return registered singleton value",
		getReturnRegisteredSingletonValue)
	t.Run("Get return registered transient instance",
		getReturnRegisteredTransientInstance)
	t.Run("Get return unique registered transient instance",
		getReturnUniqueRegisteredTransientInstance)
	t.Run("Get return unique registered transient with resolved constructor args",
		getReturnUniqueRegisteredTransientWithResolvedConstructorArgs)
}

func getReturnRegisteredSingletonInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstInterface](newFirstImpl)
	ass.NoError(err)

	instance, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstImpl{}, instance)
}

func getReturnNotUniqueRegisteredSingletonInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstInterface](newFirstImpl)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstImpl{}, v0)

	v0.SetText(random.String(10))

	v1, err := Get[iFirstInterface]()
	ass.NoError(err)

	ass.Equal(v0.GetText(), v1.GetText())
	ass.Same(v0, v1)
}

func getReturnRegisteredSingletonPointer(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[*secondImpl](newSecondImpl)
	ass.NoError(err)

	value, err := Get[*secondImpl]()
	ass.NoError(err)
	ass.IsType(&secondImpl{}, value)
}

func getReturnRegisteredSingletonValue(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[secondImpl](newSecondImplNonPointer)
	ass.NoError(err)

	value, err := Get[secondImpl]()
	ass.NoError(err)
	ass.IsType(secondImpl{}, value)
}

func getReturnRegisteredTransientInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddTransient[iFirstInterface](newFirstImpl)
	ass.NoError(err)

	value, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(new(iFirstImpl), value)
}

func getReturnUniqueRegisteredTransientInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(2)

	err := AddTransient[iFirstInterface](newFirstImpl)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstImpl{}, v0)

	v0.SetText(random.String(10))

	v1, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstImpl{}, v1)

	ass.NotEqual(v0.GetText(), v1.GetText())
	ass.NotSame(v0, v1)
}

func getReturnUniqueRegisteredTransientWithResolvedConstructorArgs(t *testing.T) {
	ass := assert.New(t)
	clearValues(2)

	err := AddSingleton[iSecondInterface](newSecondImpl)
	ass.NoError(err)

	second, err := Get[iSecondInterface]()
	ass.NoError(err)
	secondText := random.String(10)
	second.SetSecondText(secondText)

	err = AddTransient[iFirstInterface](newFirstImplWithSecondInterface)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)

	ass.Equal(v0.GetTextFromSecond(), secondText)
}
