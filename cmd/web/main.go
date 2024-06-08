package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// create custom cli flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// add structured logger, use JSON logging, add debug log level and file/line number source
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// app struct with new logger
	app := &application{
		logger: logger,
	}

	// create mux
	mux := http.NewServeMux()

	// create file server for static files
	/*	fileServer := http.FileServer(http.Dir("./ui/static"))
		mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))*/

	//  create custom file server with dir listing disabled
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// add handler functions to endpoints
	mux.HandleFunc("GET /{$}", app.home) // restrict for / only
	mux.HandleFunc("GET /snippet/view/{id}/{$}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// start server
	logger.Info("starting server", slog.String("addr", *addr))
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
