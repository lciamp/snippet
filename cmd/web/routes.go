package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//  create custom file server with dir listing disabled
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// add handler functions to endpoints
	mux.HandleFunc("GET /{$}", app.home) // restrict for / only
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return commonHeaders(mux)
}
