package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

// 500 error at log error level
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// error when problem with request from client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// custom fileSystem to disable File Server Directory Listings
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open method for custom file system
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// get template from templateCache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// initalize new buffer
	buf := new(bytes.Buffer)

	// write template to buffer
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// write http status code
	w.WriteHeader(status)

	// write contents of buffer to response writer
	buf.WriteTo(w)

}
