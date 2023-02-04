package register

type IRegister[T any] interface {
	Register() error
}
