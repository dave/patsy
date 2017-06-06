package patsy_test

import (
	"path"
	"path/filepath"
	"testing"

	"strings"

	"github.com/dave/patsy"
	"github.com/dave/patsy/builder"
	"github.com/dave/patsy/vos"
)

func TestName2(t *testing.T) {
	env := vos.Mock()
	b, err := builder.New(env, "ns")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	_, dirA, err := b.Package("a", map[string]string{
		"a.go": "package a",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, dirB, err := b.Package("b", map[string]string{
		"b.go": "package b",
	})
	if err != nil {
		t.Fatal(err)
	}

	packagePathC, _, err := b.Package("c", map[string]string{
		"c.go": "package c",
	})
	if err != nil {
		t.Fatal(err)
	}

	// We add a vendored version of "c" inside "b" that has the name "v"
	_, _, err = b.Package(path.Join("b", "vendor", packagePathC), map[string]string{
		"c.go": "package v",
	})
	if err != nil {
		t.Fatal(err)
	}

	name, err := patsy.Name(env, packagePathC, dirA)
	if err != nil {
		t.Fatal(err)
	}
	expected := "c"
	if name != expected {
		t.Fatalf("Got %s, Expected %s", name, expected)
	}

	name, err = patsy.Name(env, packagePathC, dirB)
	if err != nil {
		t.Fatal(err)
	}
	expected = "v"
	if name != expected {
		t.Fatalf("Got %s, Expected %s", name, expected)
	}
}

func TestName(t *testing.T) {
	env := vos.Mock()
	b, err := builder.New(env, "ns")
	if err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup()

	packagePath, _, err := b.Package("a", map[string]string{
		"a.go": "package b",
	})
	if err != nil {
		t.Fatal(err)
	}
	name, err := patsy.Name(env, packagePath, "/")
	if err != nil {
		t.Fatal(err)
	}
	expected := "b"
	if name != expected {
		t.Fatalf("Got %s, Expected %s", name, expected)
	}
}

func TestPath(t *testing.T) {

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

func TestDir(t *testing.T) {

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
