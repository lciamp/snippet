package main

import (
	"net/http"
	"os"
)

// custom fileSystem to disable Fileserver Directory Listings
type neuteredFileSystem struct {
	fs http.FileSystem
}

// open method for custom file system
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}
