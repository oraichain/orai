package filecache

import (
	"github.com/peterbourgon/diskv"
)

// Cache is used to store oracle script and data source files in bytes
type Cache struct {
	fileCache *diskv.Diskv
}

// New creates and returns a new file-backed data caching instance.
func New(basePath string) Cache {
	return Cache{
		fileCache: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 32 * 1024 * 1024, // 32MB TODO: Make this configurable
		}),
	}
}

// AddFile saves the given data to a file in HOME/files directory
func (c Cache) AddFile(data []byte, filename string) {
	if !c.fileCache.Has(filename) {
		c.fileCache.Write(filename, data)
	}
}

// GetFile loads the file from the file storage. Returns error if the file does not exist.
func (c Cache) GetFile(filename string) ([]byte, error) {
	data, err := c.fileCache.Read(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MustGetFile loads the file from the file storage. Panics if the file does not exist.
func (c Cache) MustGetFile(filename string) []byte {
	data, err := c.GetFile(filename)
	if err != nil {
		panic(err)
	}
	return data
}

// EraseFile delete a file from the filesystem.
func (c Cache) EraseFile(filename string) {

	err := c.fileCache.Erase(filename)
	if err != nil {
		panic("Cannot erase file")
	}
}

// EditFile edit a given file in HOME/files directory
func (c Cache) EditFile(data []byte, filename string) {
	if c.fileCache.Has(filename) {
		c.fileCache.Write(filename, data)
	}
}
