package vos

import (
	"github.com/dave/patsy/vos/mock"
	"github.com/dave/patsy/vos/os"
)

type Env interface {
	Environ() []string
	Getenv(key string) string
	Setenv(key, value string) error
	Getwd() (string, error)
	Setwd(dir string) error
}

var _ Env = (*os.Env)(nil)
var _ Env = (*mock.Env)(nil)

func Os() Env {
	return os.New()
}

func Mock() Env {
	return mock.New()
}
