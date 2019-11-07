package storage

// Cache uses to keeping or retrieving the uploaded files from the cache storage
type Cache interface {
	// Exist checks the existence of the file
	Exist(path string) bool

	// Retrieve returns a file from the given path, an error will be returned if file not found
	Retrieve(path string) (File, error)

	// Store keeps a file in the cache storage
	Store(file File) error
}

// NewDiskCache returns new disk storage cache instance
func NewDiskCache(cachePath string) DiskCache {
	return DiskCache{cachePath: cachePath}
}

// DiskCache uses the hard-disk drive as a cache storage
type DiskCache struct {
	cachePath string
}

func (c DiskCache) Exist(path string) bool {
	return false
}

func (c DiskCache) Retrieve(path string) (File, error) {
	return File{}, nil
}

func (c DiskCache) Store(file File) error {
	return nil
}
