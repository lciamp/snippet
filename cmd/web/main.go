package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// cli flags to be determined at run time
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create servemux instance
	mux := http.NewServeMux()
	// create file server for static files
	//fileServer := http.FileServer(http.Dir("./ui/static"))

	// use mux.Handle() to register the files server as the handler for all URL paths that start with
	// "/static/". For matching paths, strip "/static" prefix before the request reaches the file server
	//mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // restrict for / only
	mux.HandleFunc("GET /snippet/view/{id}/{$}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("Starting on: localhost:%s", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
