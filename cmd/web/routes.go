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

	// add handler functions to endpoints
	// switch back to handle instead of HandleFunc for dynamic middleware (ThenFunc returns a Handler)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home)) // restrict for / only
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	// new routes for user login/out/create
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// create middleware chain using our standard middleware
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// use the standard middleware
	return standard.Then(mux)
}
