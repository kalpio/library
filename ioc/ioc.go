package ioc

import (
	"github.com/pkg/errors"
	"reflect"
)

var (
	values               = make(map[reflect.Type]interface{}, 10)
	ErrAlreadyRegistered = errors.New("ioc: already registered type")
	ErrNoRegistered      = errors.New("ioc: no service is registered for type")
	ErrInvalidType       = errors.New("ioc: ")
)

func AddSingleton[T any](value any) error {
	t := getType[T]()
	if t.Kind() == reflect.Interface {
		reflect.TypeOf(value).Implements(t)
	}
	_, exist := values[t]
	if exist {
		return errors.Wrapf(ErrAlreadyRegistered, "%s", t)
	}

	values[t] = value

	return nil
}

func Get[T any]() (T, error) {
	t := getType[T]()

	value, ok := values[t]
	if !ok {
		return *new(T), errors.Wrapf(ErrNoRegistered, "%s", t)
	}

	return value.(T), nil
}

func getType[T any]() reflect.Type {
	var obj T
	return reflect.TypeOf(&obj).Elem()
}
