package patsy

import (
	"path"
	"path/filepath"

	"sync"

	"github.com/dave/patsy/vos"
)

// NewCache returns a new *Cache, allowing cached access to patsy utility
// functions.
func NewCache(env vos.Env) *Cache {
	return &Cache{
		env:    env,
		dirsm:  new(sync.RWMutex),
		pathsm: new(sync.RWMutex),
		namesm: new(sync.RWMutex),
		dirs:   make(map[string]string),
		paths:  make(map[string]string),
		names:  make(map[namekey]string),
	}
}

type namekey struct {
	path string
	dir  string
}

// Cache supports patsy.Dir and patsy.Path, but cached so they can be used in
// tight loops without hammering the filesystem.
type Cache struct {
	env    vos.Env
	dirsm  *sync.RWMutex
	pathsm *sync.RWMutex
	namesm *sync.RWMutex
	dirs   map[string]string
	paths  map[string]string
	names  map[namekey]string
}

// Name does the same as patsy.Name but cached.
func (c *Cache) Name(packagePath, srcDir string) (string, error) {
	// check the cache first
	if n, ok := c.getName(namekey{path: packagePath, dir: srcDir}); ok {
		return n, nil
	}
	n, err := Name(c.env, packagePath, srcDir)
	if err != nil {
		return "", err
	}
	c.setName(namekey{path: packagePath, dir: srcDir}, n)
	return n, nil
}

// Path does the same as patsy.Path but cached.
func (c *Cache) Path(dir string) (string, error) {
	// check the cache first
	if ppath, ok := c.getPath(dir); ok {
		return ppath, nil
	}
	ppath, err := Path(c.env, dir)
	if err != nil {
		return "", err
	}
	c.setDir(ppath, dir)
	c.setPath(dir, ppath)
	return ppath, nil
}

// Dir does the same as patsy.Dir but cached.
func (c *Cache) Dir(ppath string) (string, error) {
	// check the cache first
	if dir, ok := c.getDir(ppath); ok {
		return dir, nil
	}
	dir, err := Dir(c.env, ppath)
	if err != nil {
		return "", err
	}
	c.setDir(ppath, dir)
	c.setPath(dir, ppath)
	return dir, nil
}

// GoName converts a full filepath to a package path and filename:
//     /Users/dave/go/src/github.com/dave/foo.go -> github.com/dave/foo.go
func (c *Cache) GoName(fpath string) (string, error) {
	fdir, fname := filepath.Split(fpath)
	ppath, err := c.Path(fdir)
	if err != nil {
		return "", err
	}
	return path.Join(ppath, fname), nil
}

// FilePath converts a package path and filename to a full filepath:
//     github.com/dave/foo.go -> /Users/dave/go/src/github.com/dave/foo.go
func (c *Cache) FilePath(gpath string) (string, error) {
	ppath, fname := path.Split(gpath)
	fdir, err := c.Dir(ppath)
	if err != nil {
		return "", err
	}
	return filepath.Join(fdir, fname), nil
}

func (c *Cache) getDir(key string) (string, bool) {
	c.dirsm.RLock()
	defer c.dirsm.RUnlock()
	v, ok := c.dirs[key]
	return v, ok
}

func (c *Cache) getPath(key string) (string, bool) {
	c.pathsm.RLock()
	defer c.pathsm.RUnlock()
	v, ok := c.paths[key]
	return v, ok
}

func (c *Cache) getName(key namekey) (string, bool) {
	c.namesm.RLock()
	defer c.namesm.RUnlock()
	v, ok := c.names[key]
	return v, ok
}

func (c *Cache) setDir(key, value string) {
	c.dirsm.Lock()
	defer c.dirsm.Unlock()
	c.dirs[key] = value
}

func (c *Cache) setPath(key, value string) {
	c.pathsm.Lock()
	defer c.pathsm.Unlock()
	c.paths[key] = value
}

func (c *Cache) setName(key namekey, value string) {
	c.namesm.Lock()
	defer c.namesm.Unlock()
	c.names[key] = value
}
