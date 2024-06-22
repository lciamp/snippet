package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippet.lciamp.xyz/internal/models"
	"strconv"
	"strings"
	"unicode/utf8"
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
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// create struct to deal with errors.
// note: all fields start with a capital letter so they can be exported to templates
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

// add a snippet handler to POST snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse the form into a form map
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// from data always returns a string, we need to convert to an integer
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// create snippetCreateForm struct containing the values from the form and empty map for errors
	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// update validation checks to use new struct

	// check title is not blank or over 100 chars
	// if either fails add error to fieldErrors
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field can not be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field can not be greater than 100 characters"
	}

	// check if the content value is blank
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field can not be blank"
	}

	// check if the expires field matches (1, 7 oe 365)
	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// if any errors redisplay then create template with new form
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// pass data to insert
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// redirect to new snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
