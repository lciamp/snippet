package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// create file server for static files using custom file system from helpers.go
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // restrict for / only
	mux.HandleFunc("GET /snippet/view/{id}/{$}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting on: localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
