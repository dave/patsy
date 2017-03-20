package patsy_test

import (
	"path/filepath"
	"testing"

	"strings"

	"github.com/dave/patsy/builder"
	"github.com/dave/patsy/vos"
	"github.com/dave/patsy"
)

func TestGetPackageFromDir(t *testing.T) {

	env := vos.Mock()
	b, err := builder.New(env, "ns")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, packageDir, err := b.Package("a", nil)
	if err != nil {
		t.Fatal(err)
	}

	calculatedPath, err := patsy.Path(env, packageDir)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedPath != packagePath {
		t.Fatalf("Got %s, Expected %s", calculatedPath, packagePath)
	}

	env.Setenv("GOPATH", "/foo/"+string(filepath.ListSeparator)+env.Getenv("GOPATH"))

	calculatedPath, err = patsy.Path(env, packageDir)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedPath != packagePath {
		t.Fatalf("Got %s, Expected %s", calculatedPath, packagePath)
	}

	env.Setenv("GOPATH", "/bar/")
	_, err = patsy.Path(env, packageDir)
	if err == nil {
		t.Fatal("Expected error, got none.")
	} else if !strings.HasPrefix(err.Error(), "Package not found") {
		t.Fatalf("Expected 'Package not found', got '%s'", err.Error())
	}
}

func TestGetDirFromPackage(t *testing.T) {

	env := vos.Mock()
	b, err := builder.New(env, "ns")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, packageDir, err := b.Package("a", nil)
	if err != nil {
		t.Fatal(err)
	}

	calculatedDir, err := patsy.Dir(env, packagePath)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedDir != packageDir {
		t.Fatalf("Got %s, expected %s", calculatedDir, packageDir)
	}

	err = b.File("a", "a.go", "package a")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: somehow ensure exe.CombinedOutput() succeeds?
	calculatedDir, err = patsy.Dir(env, packagePath)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedDir != packageDir {
		t.Fatalf("Got %s, expected %s", calculatedDir, packageDir)
	}

}
