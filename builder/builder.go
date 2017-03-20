// package builder can be used in testing to create a temporary gopath, src, 
// namespace and package directory, and populate it with source files.
package builder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"path"

	"github.com/dave/patsy/vos"
	"github.com/pkg/errors"
)

// New creates a new gopath in the system temporary location, creates the src
// dir and the namespace dir. The gopath is appended to the beginning of the 
// existing gopath, so existing imports will still work. Remember to defer the 
// Cleanup() method to delete the temporary files.
func New(env vos.Env, namespace string) (*Builder, error) {

	gopath, err := ioutil.TempDir("", "go")
	if err != nil {
		return nil, errors.Wrap(err, "Error creating temporary gopath root dir")
	}

	b := &Builder{
		env:       env,
		root:      gopath,
		namespace: namespace,
	}

	if err := os.Mkdir(filepath.Join(gopath, "src"), os.FileMode(0777)); err != nil {
		b.Cleanup()
		return nil, errors.Wrap(err, "Error creating temporary gopath src dir")
	}
	b.env.Setenv("GOPATH", gopath+string(filepath.ListSeparator)+b.env.Getenv("GOPATH"))

	if err := os.MkdirAll(filepath.Join(gopath, "src", namespace), os.FileMode(0777)); err != nil {
		b.Cleanup()
		return nil, errors.Wrap(err, "Error creating temporary namespace dir")
	}

	return b, nil
}

// Builder can be used in testing to create a temporary gopath, src, namespace 
// and package directory, and populate it with source files.
type Builder struct {
	env       vos.Env // mockable environment
	root      string  // temporary gopath root dir
	namespace string  // temporary namespace
}

// File creates a new source file in the package.
func (b *Builder) File(packageName, filename, contents string) error {
	dir := filepath.Join(b.root, "src", b.namespace, packageName)
	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		// most editors will indent multi line strings in Go source with
		// tabs, so we convert to spaces for yaml files.
		contents = strings.Replace(contents, "\t", "    ", -1)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, filename), []byte(contents), 0777); err != nil {
		return errors.Wrapf(err, "Error creating temporary source file %s", filename)
	}
	return nil
}

// Package creates a new package and populates with source files.
func (b *Builder) Package(packageName string, files map[string]string) (packagePath string, packageDir string, err error) {

	dir := filepath.Join(b.root, "src", b.namespace, packageName)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", "", errors.Wrap(err, "Error creating temporary package dir")
	}

	if files != nil {
		for filename, contents := range files {
			if err := b.File(packageName, filename, contents); err != nil {
				return "", "", err
			}
		}
	}

	return path.Join(b.namespace, packageName), dir, nil
}

// Cleanup deletes all temporary files.
func (b *Builder) Cleanup() {
	os.RemoveAll(b.root)
}
