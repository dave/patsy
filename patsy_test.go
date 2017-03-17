package patsy_test

import (
	"path/filepath"
	"testing"

	"strings"

	"github.com/dave/patsy"
	"github.com/dave/patsy/builder"
	"github.com/dave/patsy/vos"
)

func TestGetCurrentGopath(t *testing.T) {
	env := vos.Mock()
	abc := filepath.Join("a", "b", "c")
	def := filepath.Join("d", "e", "f")
	env.Setenv("GOPATH", abc+string(filepath.ListSeparator)+def)
	gop := patsy.GetCurrentGopath(env)

	if gop != abc {
		t.Fatalf("Expected %s", abc)
	}
	env.Setwd(filepath.Join("d", "e", "f", "g", "h"))

	gop = patsy.GetCurrentGopath(env)
	if gop != def {
		t.Fatalf("Expected %s", def)
	}
}

func TestGetPackageFromDir(t *testing.T) {

	env := vos.Mock()
	b, err := builder.New(env)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, packageDir, err := b.Package("a", nil)
	if err != nil {
		t.Fatal(err)
	}

	calculatedPath, err := patsy.GetPackageFromDir(env, packageDir)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedPath != packagePath {
		t.Fatalf("Got %s, Expected %s", calculatedPath, packagePath)
	}

	env.Setenv("GOPATH", "/foo/"+string(filepath.ListSeparator)+env.Getenv("GOPATH"))

	calculatedPath, err = patsy.GetPackageFromDir(env, packageDir)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedPath != packagePath {
		t.Fatalf("Got %s, Expected %s", calculatedPath, packagePath)
	}

	env.Setenv("GOPATH", "/bar/")
	_, err = patsy.GetPackageFromDir(env, packageDir)
	if err == nil {
		t.Fatal("Expected error, got none.")
	} else if !strings.HasPrefix(err.Error(), "Package not found") {
		t.Fatalf("Expected 'Package not found', got '%s'", err.Error())
	}
}

func TestGetDirFromEmptyPackage(t *testing.T) {

	env := vos.Mock()
	b, err := builder.New(env)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, packageDir, err := b.Package("a", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = patsy.GetPackageFromDir(env, "/foo")
	if err == nil {
		t.Fatal("Expected error, got none.")
	} else if !strings.HasPrefix(err.Error(), "Package not found for /foo") {
		t.Fatalf("Expected 'Package not found for /foo', got '%s'", err.Error())
	}

	calculatedDir, err := patsy.GetDirFromEmptyPackage(env, packagePath)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedDir != packageDir {
		t.Fatalf("Got %s, expected %s", calculatedDir, packageDir)
	}

	env.Setenv("GOPATH", "/foo/"+string(filepath.ListSeparator)+env.Getenv("GOPATH"))

	// This will now need two loops around to get the package
	calculatedDir, err = patsy.GetDirFromEmptyPackage(env, packagePath)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedDir != packageDir {
		t.Fatalf("Got %s, expected %s", calculatedDir, packageDir)
	}

}

func TestGetDirFromPackage(t *testing.T) {

	env := vos.Mock()
	b, err := builder.New(env)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, packageDir, err := b.Package("a", nil)
	if err != nil {
		t.Fatal(err)
	}

	calculatedDir, err := patsy.GetDirFromPackage(env, packagePath)
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
	calculatedDir, err = patsy.GetDirFromPackage(env, packagePath)
	if err != nil {
		t.Fatal(err)
	}
	if calculatedDir != packageDir {
		t.Fatalf("Got %s, expected %s", calculatedDir, packageDir)
	}

}
