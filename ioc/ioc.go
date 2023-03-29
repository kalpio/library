package ioc

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

var (
	values          = make(map[reflect.Type]*scopeAndInterface, 10)
	singletons      = make(map[reflect.Type]reflect.Value, 10)
	ErrNoRegistered = errors.New("ioc: no service is registered for type")
	ErrInvalidType  = errors.New("ioc: invalid type")
	valuesLock      = &sync.RWMutex{}
)

type scopeAndInterface struct {
	scope scope
	value interface{}
	args  []reflect.Type
}

func newScopeAndInterface(scope scope, value interface{}, args []reflect.Type) *scopeAndInterface {
	return &scopeAndInterface{scope: scope, value: value, args: args}
}

type scope int32

const (
	Singleton scope = iota
	Transient
)

func AddSingleton[T interface{}](ctor any) error {
	constructor := reflect.ValueOf(ctor)

	if constructor.Kind() != reflect.Func {
		return errors.New("ioc: it is not a function")
	}

	t := getType[T]()
	if err := isValidReturnType[T](constructor); err != nil {
		return err
	}

	args := getMethodArgumentTypes(constructor)

	valuesLock.Lock()
	values[t] = newScopeAndInterface(Singleton, constructor, args)
	valuesLock.Unlock()

	return nil
}

func AddTransient[T any](ctor any) error {
	constructor := reflect.ValueOf(ctor)

	if constructor.Kind() != reflect.Func {
		return errors.New("ioc: it is not a function")
	}

	t := getType[T]()
	if err := isValidReturnType[T](constructor); err != nil {
		return err
	}

	args := getMethodArgumentTypes(constructor)

	valuesLock.Lock()
	values[t] = newScopeAndInterface(Transient, constructor, args)
	valuesLock.Unlock()

	return nil
}

func isValidReturnType[T any](constructor reflect.Value) error {
	expected := getType[T]()
	returned := constructor.Type().Out(0)
	isInterface := expected.Kind() == reflect.Interface

	if returned == expected {
		return nil
	}

	if isInterface && returned.Implements(expected) {
		return nil
	}

	if isInterface {
		return errors.Wrap(ErrInvalidType, fmt.Sprintf("return type must be %s or must implements %s", expected.String(), expected.String()))
	}

	return errors.Wrap(ErrInvalidType, fmt.Sprintf("ioc: return type must be %s", expected.String()))
}

func getMethodArgumentTypes(constructor reflect.Value) []reflect.Type {
	m := constructor.Type()
	result := make([]reflect.Type, m.NumIn())

	for i := 0; i < m.NumIn(); i++ {
		result[i] = m.In(i)
	}

	return result
}

func Get[T any]() (T, error) {
	t := getType[T]()

	v, err := resolveType(t)
	if v == getNilValue() {
		return getNil().(T), errors.Errorf("ioc: could not resolve type %s", t.String())
	}

	return v.Interface().(T), err
}

func resolveType(t reflect.Type) (reflect.Value, error) {
	valuesLock.RLock()
	value, ok := values[t]
	if !ok {
		return getNilValue(), errors.Wrapf(ErrNoRegistered, "%v", t)
	}
	valuesLock.RUnlock()

	var (
		args []reflect.Value
		err  error
	)
	if len(value.args) > 0 {
		args, err = resolveTypes(value.args)

		if err != nil {
			return getNilValue(), err
		}
	}

	if value.scope == Singleton {
		valuesLock.Lock()
		_, ok := singletons[t]
		if !ok {
			v := value.value.(reflect.Value)
			res := v.Call(args)
			singletons[t] = res[0]
		}
		valuesLock.Unlock()

		return reflect.ValueOf(singletons[t].Interface()), nil
	}

	if value.scope == Transient {
		v := value.value.(reflect.Value)
		res := v.Call(args)

		return reflect.ValueOf(res[0].Interface()), nil
	}

	return getNilValue(), errors.New("ioc: unknown scope")
}

func resolveTypes(args []reflect.Type) ([]reflect.Value, error) {
	result := make([]reflect.Value, len(args))

	for i, arg := range args {
		v, err := resolveType(arg)

		if err != nil {
			return nil, err
		}

		result[i] = v
	}

	return result, nil
}

func getType[T interface{}]() reflect.Type {
	var obj T
	return reflect.TypeOf(&obj).Elem()
}

func getNil() interface{} {
	val := reflect.TypeOf((*interface{})(nil)).Elem()
	return reflect.Zero(val).Interface()
}

func getNilValue() reflect.Value {
	return reflect.ValueOf((*interface{})(nil)).Elem()
}
