package main

import (
	"github.com/justinas/alice"
	"net/http"
)

// update return to http.Header
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	//  create custom file server with dir listing disabled
	// don't wrap static routes with session middleware
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// new middleware chain for session management / dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// dynamic routes (unprotected)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // restrict for / only
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// middleware for protected routes
	protected := dynamic.Append(app.requireAuthentication)

	// protected routes (protected by session)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// create middleware chain using our standard middleware
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	// use the standard middleware on all routes
	return standard.Then(mux)
}
