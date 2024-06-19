package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippet.lciamp.xyz/internal/models"
	"strconv"
)

// home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//get last 10 snippets
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
	}

	// call newTemplateData
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// use new render helper
	app.render(w, r, http.StatusOK, "home.tmpl", data)

}

// add snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// get wildcard, check if it's a positive integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// response, use SnippetModel's Get method
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// call newTemplateData helper
	data := app.newTemplateData(r)
	data.Snippet = snippet

	// use new render helper
	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

// add a snippet handler function to GET snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display form for creating a new snippet..."))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// add a snippet handler to POST snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// test
	title := "O snail"
	content := "O snail\nClimb Mount Fuji\nBut slowly, slowly\n\n-Kobayashi Issa"
	expires := 7

	// pass data to insert
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// redirect to new snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
