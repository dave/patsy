package patsy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dave/patsy/vos"
	"github.com/pkg/errors"
)

func GetDirFromPackage(env vos.Env, packagePath string) (string, error) {

	exe := exec.Command("go", "list", "-f", "{{.Dir}}", packagePath)
	exe.Env = env.Environ()
	out, err := exe.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	dir, err := GetDirFromEmptyPackage(env, packagePath)
	if err != nil {
		return "", err
	}
	return dir, nil

}

func GetDirFromEmptyPackage(env vos.Env, packagePath string) (string, error) {
	for _, gopath := range filepath.SplitList(env.Getenv("GOPATH")) {
		dir := filepath.Join(gopath, "src", packagePath)
		if s, err := os.Stat(dir); err == nil && s.IsDir() {
			return dir, nil
		}
	}
	return "", errors.Errorf("%s not found", packagePath)
}

func GetPackageFromDir(env vos.Env, packageDir string) (string, error) {
	var err error
	for _, gopath := range filepath.SplitList(env.Getenv("GOPATH")) {
		if strings.HasPrefix(packageDir, gopath) {
			src := fmt.Sprintf("%s/src", gopath)
			rel, inner := filepath.Rel(src, packageDir)
			if inner != nil {
				// I don't think we can trigger this error if dir starts with
				// gopath
				err = inner
				continue
			}
			if rel == "" {
				// I don't think we can trigger this either
				continue
			}
			// Remember we're returning a package path, which uses forward
			// slashes even on windows
			return filepath.ToSlash(rel), nil
		}
	}
	if err != nil {
		return "", err
	}
	return "", errors.Errorf("Package not found for %s", packageDir)
}

func GetCurrentGopath(env vos.Env) string {
	gopaths := filepath.SplitList(env.Getenv("GOPATH"))
	currentDir, err := env.Getwd()
	if err != nil {
		// can't find the current working dir
		return gopaths[0]
	}
	for _, gopath := range gopaths {
		if strings.HasPrefix(currentDir, gopath) {
			return gopath
		}
	}
	return gopaths[0]
}
