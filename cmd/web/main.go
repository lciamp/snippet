package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// create custom cli flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create mux
	mux := http.NewServeMux()

	// create file server for static files
	/*	fileServer := http.FileServer(http.Dir("./ui/static"))
		mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))*/

	//  create custom file server
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // restrict for / only
	mux.HandleFunc("GET /snippet/view/{id}/{$}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("starting on: localhost%s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
