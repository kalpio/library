package ioc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type iFakeInterface interface {
	DoSth()
}

type iFakeStruct struct {
}

func (*iFakeStruct) DoSth() {
}

type secondFakeStruct struct {
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
	err := AddSingleton[iFakeInterface](&iFakeStruct{})

	ass.NoError(err)
}

func addSingletonNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddSingleton[iFakeInterface](new(iFakeStruct))

	ass.NoError(err)
	ass.Len(values, 1)
}

func addSingletonFailWhenRegisterSameInterfaceSecondTimes(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](new(iFakeStruct))
	ass.NoError(err)

	err = AddSingleton[iFakeInterface](new(iFakeStruct))
	ass.Error(err)
	ass.ErrorIs(err, ErrAlreadyRegistered)
	var obj iFakeInterface
	ass.ErrorContains(err, reflect.TypeOf(&obj).Elem().String())
}

func addSingletonFailWhenTryingAddTypeWhichNotImplementsInterface(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](&secondFakeStruct{})
	ass.ErrorIs(err, ErrInvalidType)
}

func addSingletonFailWhenTryingAddTypeWhichIsDifferent(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeStruct](&secondFakeStruct{})
	ass.ErrorIs(err, ErrInvalidType)
}

func TestGet(t *testing.T) {
	t.Run("Get return registered instance",
		getReturnRegisteredInstance)
	t.Run("Get return registered pointer",
		getReturnRegisteredPointer)
	t.Run("Get return registered value",
		getReturnRegisteredValue)
}

func getReturnRegisteredInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](new(iFakeStruct))
	ass.NoError(err)

	instance, err := Get[iFakeInterface]()
	ass.NoError(err)
	ass.IsType(&iFakeStruct{}, instance)
}

func getReturnRegisteredPointer(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[*secondFakeStruct](&secondFakeStruct{})
	ass.NoError(err)

	value, err := Get[*secondFakeStruct]()
	ass.NoError(err)
	ass.IsType(&secondFakeStruct{}, value)
}

func getReturnRegisteredValue(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[secondFakeStruct](secondFakeStruct{})
	ass.NoError(err)

	value, err := Get[secondFakeStruct]()
	ass.NoError(err)
	ass.IsType(secondFakeStruct{}, value)
}

func clearValues(length int) {
	values = make(map[reflect.Type]interface{}, length)
}
