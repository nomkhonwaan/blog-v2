// +build linux darwin

package storage

import "golang.org/x/sys/unix"

// NewDiskCache returns new disk storage cache instance
func NewDiskCache(cacheFilesPath string) (DiskCache, error) {
	if err := unix.Access(cacheFilesPath, unix.W_OK); err != nil {
		return DiskCache{}, err
	}

	return DiskCache{cacheFilesPath: cacheFilesPath}, nil
}
