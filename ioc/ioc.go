package ioc

import (
	"github.com/pkg/errors"
	"reflect"
)

var values = map[reflect.Type]interface{}{}

func Register[T any](value any) error {
	var obj T
	tType := reflect.TypeOf(obj)

	_, exist := values[tType]
	if exist {
		return errors.Errorf("")
	}

	values[tType] = value

	return nil
}

func Get[T any]() (T, error) {
	var obj T
	tType := reflect.TypeOf(obj)

	value, ok := values[tType]
	if !ok {
		return *new(T), errors.Errorf("")
	}

	return value.(T), nil
}
