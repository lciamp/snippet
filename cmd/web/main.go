package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippet.lciamp.xyz/internal/models"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// create custom cli flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	// mysql data source name command line flag
	dsn := flag.String("dsn", "web:hudson@/snippetbox?parseTime=true", "MySQL datasource name")
	flag.Parse()

	// add structured logger, use JSON logging, add debug log level and file/line number source
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// db connection pool setup
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// create template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// create app struct for dependency injections
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// start server
	logger.Info("starting server", slog.String("addr", *addr))
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// wrap sql.Open() and return sql.DB connection pool for DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
