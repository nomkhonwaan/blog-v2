package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Cache uses to keeping or retrieving the uploaded files from the cache storage
type Cache interface {
	// Exist checks the existence of the file
	Exist(path string) bool

	// Retrieve returns a file from the given path, an error will be returned if file not found
	Retrieve(path string) (File, error)

	// Store keeps a file in the cache storage
	Store(file File) error
}

// DiskCache uses the hard-disk drive as a cache storage
type DiskCache struct {
	cacheFilesPath string
}

func (c DiskCache) Exist(path string) bool {
	_, err := os.Stat(c.cacheFilesPath + string(filepath.Separator) + path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c DiskCache) Retrieve(path string) (File, error) {
	f, err := os.Open(c.cacheFilesPath + string(filepath.Separator) + path)
	if err != nil {
		return File{}, err
	}

	body, _ := ioutil.ReadAll(f)

	return File{
		Path:     path,
		FileName: f.Name(),
		Body:     body,
	}, nil
}

func (c DiskCache) Store(file File) error {
	dir := filepath.Dir(c.cacheFilesPath + string(filepath.Separator) + file.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(c.cacheFilesPath+string(filepath.Separator)+file.Path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(file.Body)
	return err
}
