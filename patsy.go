// package patsy is a package helper utility. It allows the coonversion of go
// package paths to filesystem directories and vice versa.
package patsy

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dave/patsy/vos"
	"github.com/pkg/errors"
)

func Dir(env vos.Env, packagePath string) (string, error) {

	exe := exec.Command("go", "list", "-f", "{{.Dir}}", packagePath)
	exe.Env = env.Environ()
	out, err := exe.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	// The go list command will throw an error if the package directory is
	// empty. In this case we need to explore the filesystem to see if there is
	// a directory in <gopath>/src/<package-path>. Remember there can be
	// several gopaths. We return the first matching directory.
	for _, gopath := range filepath.SplitList(env.Getenv("GOPATH")) {
		dir := filepath.Join(gopath, "src", packagePath)
		if s, err := os.Stat(dir); err == nil && s.IsDir() {
			return dir, nil
		}
	}

	return "", errors.Errorf("Dir not found for %s", packagePath)

}

func Path(env vos.Env, packageDir string) (string, error) {
	packageDir = filepath.Clean(packageDir)
	for _, gopath := range filepath.SplitList(env.Getenv("GOPATH")) {
		if strings.HasPrefix(packageDir, gopath) {
			rel, inner := filepath.Rel(filepath.Join(gopath, "src"), packageDir)
			if inner == nil && rel != "" {
				// Remember we're returning a package path, which uses forward
				// slashes even on windows
				return filepath.ToSlash(rel), nil
			}
		}
	}
	return "", errors.Errorf("Package not found for %s", packageDir)
}
