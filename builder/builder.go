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

func New(env vos.Env) (*Builder, error) {
	b := &Builder{
		env: env,
	}
	gopath, err := ioutil.TempDir("", "go")
	if err != nil {
		return nil, errors.Wrap(err, "Error creating temporary gopath root dir")
	}
	b.root = gopath
	if err := os.Mkdir(filepath.Join(gopath, "src"), os.FileMode(0777)); err != nil {
		b.Cleanup()
		return nil, errors.Wrap(err, "Error creating temporary gopath src dir")
	}
	b.env.Setenv("GOPATH", gopath+string(filepath.ListSeparator)+b.env.Getenv("GOPATH"))

	dir, err := ioutil.TempDir(filepath.Join(gopath, "src"), "ns")
	if err != nil {
		b.Cleanup()
		return nil, errors.Wrap(err, "Error creating temporary namespace dir")
	}

	// namespace is of the form <gopath>/src/<namespace>
	b.namespace = dir[strings.LastIndex(dir, string(os.PathSeparator))+1:]

	return b, nil
}

type Builder struct {
	env       vos.Env // mockable environment
	root      string  // temporary gopath root dir
	namespace string  // temporary namespace
}

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

func (b *Builder) Cleanup() {
	os.RemoveAll(b.root)
}
