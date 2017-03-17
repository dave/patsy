package mock

import (
	"fmt"
	"os"
	"strings"
)

func New() *Env {
	return &Env{
		EnvironmentVariables: make(map[string]string),
	}
}

type Env struct {
	EnvironmentVariables map[string]string
	WorkingDirectory     string
}

func (m *Env) Getenv(key string) string {
	if m.EnvironmentVariables == nil {
		return os.Getenv(key)
	}
	e, ok := m.EnvironmentVariables[key]
	if !ok {
		return os.Getenv(key)
	}
	return e
}

func (m *Env) Setenv(key, value string) error {
	m.EnvironmentVariables[key] = value
	return nil
}

func (m *Env) Getwd() (string, error) {
	if m.WorkingDirectory == "" {
		return os.Getwd()
	}
	return m.WorkingDirectory, nil
}

func (m *Env) Setwd(dir string) error {
	m.WorkingDirectory = dir
	return nil
}

// Environ returns a copy of strings representing the environment, in the form "key=value".
func (m *Env) Environ() []string {
	if m.EnvironmentVariables == nil {
		return os.Environ()
	}
	var out []string
	merged := make(map[string]string)
	for _, e := range os.Environ() {
		// Add the environment variables from the system
		parts := strings.Split(e, "=")
		merged[parts[0]] = parts[1]
	}
	for k, v := range m.EnvironmentVariables {
		// Overwrite with the mocked environment variables
		merged[k] = v
	}
	for k, v := range merged {
		// Join them back together in Environ syntax
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}
