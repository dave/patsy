// Package vos is a virtual os tool. It allows mocking of the os.Environ,
// os.Getenv and os.Getwd functions.
package vos

import (
	"github.com/dave/patsy/vos/mock"
	"github.com/dave/patsy/vos/os"
)

// Env provides an interface with methods similar to os.Environ, os.Getenv and
// os.Getwd functions.
type Env interface {
	Environ() []string
	Getenv(key string) string
	Setenv(key, value string) error
	Getwd() (string, error)
	Setwd(dir string) error
}

var _ Env = (*os.Env)(nil)
var _ Env = (*mock.Env)(nil)

// Os returns an Env that provides a direct pass-through to the os package. Use
// this in production.
func Os() Env {
	return os.New()
}

// Mock returns an Env that provides a mock for the os package. Use this in
// testing.
func Mock() Env {
	return mock.New()
}
