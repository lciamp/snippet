package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"snippet.lciamp.xyz/internal/models"
	"time"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// initialize a decoder instance
	formDecoder := form.NewDecoder()

	//create new session manager and configure mysql db to as the session store wil 12H ttl.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// make session cookie secure
	sessionManager.Cookie.Secure = true

	// create app struct for dependency injections
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// create new htt.Server struct. Set addr and Handler fields
	svr := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// create a *log.Logger from our slog handler, which writes error level logs and assign it to ErrorLog.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", svr.Addr)

	// use ListenAndServeTLS() to start https server
	err = svr.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
