package ioc

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type iFakeInterface interface {
	DoSth()
}

type iFakeStruct struct {
}

func (f *iFakeStruct) DoSth() {
}

func TestAddSingletonNotReturnError(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddSingleton[iFakeInterface](new(iFakeStruct))

	ass.NoError(err)
}

func TestAddSingletonNotReturnErrorAndHasRegisteredService(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)
	err := AddSingleton[iFakeInterface](new(iFakeStruct))

	ass.NoError(err)
	ass.Len(values, 1)
}

func TestAddSingletonFailWhenRegisterSameInterfaceSecondTimes(t *testing.T) {
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

func TestAddSingletonFailWhenTryingAddDifferentType(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](new(iFakeStruct))
	ass.ErrorIs(err, Err)

}

func TestGetReturnRegisteredInstance(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](new(iFakeStruct))
	ass.NoError(err)

	instance, err := Get[iFakeInterface]()
	ass.NoError(err)
	ass.IsType(&iFakeStruct{}, instance)
}

func TestGetReturnRegisteredValue(t *testing.T) {
	ass := assert.New(t)
	clearValues(1)

	err := AddSingleton[iFakeInterface](iFakeStruct{})
	ass.NoError(err)

	value, err := Get[iFakeInterface]()
	ass.NoError(err)
	ass.IsType(iFakeStruct{}, value)
}

func clearValues(length int) {
	values = make(map[reflect.Type]interface{}, length)
}
