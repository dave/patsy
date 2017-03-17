package os

import "os"

func New() *Env {
	return &Env{}
}

type Env struct{}

func (*Env) Getenv(key string) string {
	return os.Getenv(key)
}

func (*Env) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (*Env) Getwd() (string, error) {
	return os.Getwd()
}

func (*Env) Setwd(dir string) error {
	return os.Chdir(dir)
}

// Environ returns a copy of strings representing the environment, in the form "key=value".
func (*Env) Environ() []string {
	return os.Environ()
}
