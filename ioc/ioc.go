package ioc

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	values               = make(map[reflect.Type]interface{}, 10)
	ErrAlreadyRegistered = errors.New("ioc: already registered type")
	ErrNoRegistered      = errors.New("ioc: no service is registered for type")
	ErrInvalidType       = errors.New("ioc: invalid type")
)

func AddSingleton[T any](value any) error {
	t := getType[T]()

	if t.Kind() == reflect.Interface {
		return addSingletonForInterface(t, value)
	}

	return addSingletonForNoneInterface(t, value)
}

func addSingletonForInterface(t reflect.Type, value any) error {
	valueType := reflect.TypeOf(value)
	if !valueType.Implements(t) {
		return errors.Wrapf(ErrInvalidType, "%v should implements %v", valueType, t)
	}

	return addSingletonInternal(t, value)
}

func addSingletonForNoneInterface(t reflect.Type, value any) error {
	valueType := reflect.TypeOf(value)
	if valueType != t {
		return errors.Wrapf(ErrInvalidType, "%v should be same type as %v", valueType, t)
	}

	return addSingletonInternal(t, value)
}

func addSingletonInternal(t reflect.Type, value any) error {
	_, exists := values[t]
	if exists {
		return errors.Wrapf(ErrAlreadyRegistered, "%v", t)
	}

	values[t] = value

	return nil
}

func RemoveSingleton[T any]() {
	t := getType[T]()

	_, exists := values[t]
	if exists {
		delete(values, t)
	}
}

func Get[T any]() (T, error) {
	t := getType[T]()

	value, ok := values[t]
	if !ok {
		return *new(T), errors.Wrapf(ErrNoRegistered, "%v", t)
	}

	return value.(T), nil
}

func getType[T any]() reflect.Type {
	var obj T
	return reflect.TypeOf(&obj).Elem()
}
