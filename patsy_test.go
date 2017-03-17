package patsy

import (
	"path/filepath"
	"testing"

	"github.com/dave/patsy/vos"
)

func TestGetCurrentGopath(t *testing.T) {
	env := vos.Mock()
	abc := filepath.Join("a", "b", "c")
	def := filepath.Join("d", "e", "f")
	env.Setenv("GOPATH", abc+string(filepath.ListSeparator)+def)
	gop := GetCurrentGopath(env)

	if gop != abc {
		t.Fatalf("Expected %s", abc)
	}
	env.Setwd(filepath.Join("d", "e", "f", "g", "h"))

	gop = GetCurrentGopath(env)
	if gop != def {
		t.Fatalf("Expected %s", def)
	}
}

/*
func TestGetPackageFromDir(t *testing.T) {

	cb := tests.New().TempGopath(false)
	defer cb.Cleanup()
	packagePath, packageDir := cb.TempPackage("a", map[string]string{})

	calculatedPath, err := GetPackageFromDir(cb.Ctx(), packageDir)
	require.NoError(t, err)
	assert.Equal(t, packagePath, calculatedPath)

	vos := vosctx.FromContext(cb.Ctx())
	cb.OsVar("GOPATH", "/fdskljsfdash/"+string(filepath.ListSeparator)+vos.Getenv("GOPATH"))

	calculatedPath, err = GetPackageFromDir(cb.Ctx(), packageDir)
	require.NoError(t, err)
	assert.Equal(t, packagePath, calculatedPath)

	cb.OsVar("GOPATH", "/fdskljsfdash/")
	_, err = GetPackageFromDir(cb.Ctx(), packageDir)
	assert.IsError(t, err, "CXOETFPTGM")
}
/*
func TestGetDirFromEmptyPackage(t *testing.T) {
	cb := tests.New().TempGopath(false)
	defer cb.Cleanup()
	packagePath, packageDir := cb.TempPackage("a", map[string]string{})

	_, err := GetDirFromEmptyPackage(cb.Ctx(), "a.b/c")
	assert.IsError(t, err, "SUTCWEVRXS")

	calculatedDir, err := GetDirFromEmptyPackage(cb.Ctx(), packagePath)
	require.NoError(t, err)
	assert.Equal(t, packageDir, calculatedDir)

	vos := vosctx.FromContext(cb.Ctx())
	cb.OsVar("GOPATH", "/fdskljsfdash/"+string(filepath.ListSeparator)+vos.Getenv("GOPATH"))

	// This will now need two loops around to get the package
	calculatedDir, err = GetDirFromEmptyPackage(cb.Ctx(), packagePath)
	require.NoError(t, err)
	assert.Equal(t, packageDir, calculatedDir)

}
func TestGetDirFromPackage(t *testing.T) {
	cb := tests.New().TempGopath(false)
	defer cb.Cleanup()
	packagePath, packageDir := cb.TempPackage("a", map[string]string{})
	calculatedDir, err := GetDirFromPackage(cb.Ctx(), packagePath)
	require.NoError(t, err)
	assert.Equal(t, packageDir, calculatedDir)

	cb.TempFile("a.go", "package a")

	calculatedDir, err = GetDirFromPackage(cb.Ctx(), packagePath)
	require.NoError(t, err)
	assert.Equal(t, packageDir, calculatedDir)

}
*/
