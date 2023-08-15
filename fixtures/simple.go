package test

import (
	"context"
	"io"
)

type MyService interface {
	io.ReadCloser
	SayHi(name string) (string, error)
	SayBye(name *string) io.Reader
}

type TypeParamService[T comparable] interface {
	DoSomething(T) error
	ReadToType(ctx context.Context, r io.Reader) []T
}
