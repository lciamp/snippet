package main

import (
	"fmt"
	"net/http"
)

// middleware for common headers - be careful with CSP Headers, check browser logs
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		next.ServeHTTP(w, r)
	})
}

// log requests middleware implement on application
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

// middleware function to recover from panics
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// use built-in recover function to see if there was a panic
			if err := recover(); err != nil {
				// set a connection close in header response
				w.Header().Set("Connection", "close")
				// send a 500 error
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requireAuthentication redirect to login if not authenticated
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if not auth redirect
		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")

		//call next
		next.ServeHTTP(w, r)
	})
}
