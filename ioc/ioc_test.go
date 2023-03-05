package ioc

import (
	"library/random"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type iFirstInterface interface {
	SetText(s string)
	GetText() string
	GetTextFromSecond() string
}

type iSecondInterface interface {
	GetSecondText() string
	SetSecondText(s string)
}

type iFirstStruct struct {
	Second iSecondInterface
	text   string
}

func newFirstStruct() *iFirstStruct {
	return new(iFirstStruct)
}

func newFirstStructWithSecondInterface(secondInterface iSecondInterface) *iFirstStruct {
	return &iFirstStruct{Second: secondInterface}
}

func (f *iFirstStruct) SetText(s string) {
	f.text = s
}

func (f *iFirstStruct) GetText() string {
	return f.text
}

func (f *iFirstStruct) GetTextFromSecond() string {
	return f.Second.GetSecondText()
}

type secondStruct struct {
	text string
}

func newSecondStruct() *secondStruct {
	return new(secondStruct)
}

func newSecondStructNonPointer() secondStruct {
	return secondStruct{}
}

func (second *secondStruct) SetSecondText(s string) {
	second.text = s
}

func (second *secondStruct) GetSecondText() string {
	return second.text
}

func TestAddTransient(t *testing.T) {
	t.Run("Add transient not return error",
		addTransientNotReturnError)
	t.Run("Add transient not return error and has registered service",
		addTransientNotReturnErrorAndHasRegisteredService)
}

func addTransientNotReturnError(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddTransient[iFirstInterface](newFirstStruct)
	ass.NoError(err)
}

func addTransientNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddTransient[iFirstInterface](newFirstStruct)

	ass.NoError(err)
	ass.Len(values, 1)
}

func TestAddSingleton(t *testing.T) {
	t.Run("Add singleton not return error",
		addSingletonNotReturnError)
	t.Run("Add singleton no return error and has registered service",
		addSingletonNotReturnErrorAndHasRegisteredService)
	t.Run("Add singleton fail when register same interface second time",
		addSingletonFailWhenRegisterSameInterfaceSecondTimes)
	t.Run("Add singleton fail when trying add type which not implement interface",
		addSingletonFailWhenTryingAddTypeWhichNotImplementsInterface)
	t.Run("Add singleton fail when trying add type which is different",
		addSingletonFailWhenTryingAddTypeWhichIsDifferent)
}

func addSingletonNotReturnError(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddSingleton[iFirstInterface](newFirstStruct)

	ass.NoError(err)
}

func addSingletonNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddSingleton[iFirstInterface](newFirstStruct)

	ass.NoError(err)
	ass.Len(values, 1)
}

func addSingletonFailWhenRegisterSameInterfaceSecondTimes(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstInterface](newFirstStruct)
	ass.NoError(err)

	err = AddSingleton[iFirstInterface](newFirstStruct)
	ass.Error(err)
	ass.ErrorIs(err, ErrAlreadyRegistered)
	var obj iFirstInterface
	ass.ErrorContains(err, reflect.TypeOf(&obj).Elem().String())
}

func addSingletonFailWhenTryingAddTypeWhichNotImplementsInterface(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstInterface](newSecondStruct)
	ass.ErrorIs(err, ErrInvalidType)
}

func addSingletonFailWhenTryingAddTypeWhichIsDifferent(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstStruct](newSecondStruct)
	ass.ErrorIs(err, ErrInvalidType)
}

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

	err := AddSingleton[iFirstInterface](newFirstStruct)
	ass.NoError(err)

	instance, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstStruct{}, instance)
}

func getReturnNotUniqueRegisteredSingletonInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFirstInterface](newFirstStruct)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstStruct{}, v0)

	v0.SetText(random.String(10))

	v1, err := Get[iFirstInterface]()
	ass.NoError(err)

	ass.Equal(v0.GetText(), v1.GetText())
	ass.Same(v0, v1)
}

func getReturnRegisteredSingletonPointer(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[*secondStruct](newSecondStruct)
	ass.NoError(err)

	value, err := Get[*secondStruct]()
	ass.NoError(err)
	ass.IsType(&secondStruct{}, value)
}

func getReturnRegisteredSingletonValue(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[secondStruct](newSecondStructNonPointer)
	ass.NoError(err)

	value, err := Get[secondStruct]()
	ass.NoError(err)
	ass.IsType(secondStruct{}, value)
}

func getReturnRegisteredTransientInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddTransient[iFirstInterface](newFirstStruct)
	ass.NoError(err)

	value, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(new(iFirstStruct), value)
}

func getReturnUniqueRegisteredTransientInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(2)

	err := AddTransient[iFirstInterface](newFirstStruct)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstStruct{}, v0)

	v0.SetText(random.String(10))

	v1, err := Get[iFirstInterface]()
	ass.NoError(err)
	ass.IsType(&iFirstStruct{}, v1)

	ass.NotEqual(v0.GetText(), v1.GetText())
	ass.NotSame(v0, v1)
}

func getReturnUniqueRegisteredTransientWithResolvedConstructorArgs(t *testing.T) {
	ass := assert.New(t)
	clearValues(2)

	err := AddSingleton[iSecondInterface](newSecondStruct)
	ass.NoError(err)

	second, err := Get[iSecondInterface]()
	ass.NoError(err)
	secondText := random.String(10)
	second.SetSecondText(secondText)

	err = AddTransient[iFirstInterface](newFirstStructWithSecondInterface)
	ass.NoError(err)

	v0, err := Get[iFirstInterface]()
	ass.NoError(err)

	ass.Equal(v0.GetTextFromSecond(), secondText)
}

func clearValues(length int) {
	values = make(map[reflect.Type]*scopeAndInterface, length)
}
