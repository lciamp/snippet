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

	// create app struct for dependency injection on handlers
	app := &application{
		logger: logger,
	}

	// start server
	logger.Info("starting server", slog.String("addr", *addr))
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
