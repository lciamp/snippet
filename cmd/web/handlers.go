package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippet.lciamp.xyz/internal/models"
	"snippet.lciamp.xyz/internal/validator"
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

	// use PopString() method to get the value for the flash key in the session
	flash := app.sessionManager.PopString(r.Context(), "flash")

	// call newTemplateData helper
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// add the flash
	data.Flash = flash

	// use new render helper
	app.render(w, r, http.StatusOK, "view.tmpl", data)

}

// create struct to deal with errors.
// note: all fields start with a capital letter so they can be exported to templates
type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// add a snippet handler function to GET snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// create new form
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

// add a snippet handler to POST snippet
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// create form
	var form snippetCreateForm

	// use decode form helper
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// check title is not blank or over 100 chars
	form.CheckField(validator.NotBlank(form.Title), "title", "This field can not be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can not be more than 100 chars")

	// check if field is blank and if expires is a valid value
	form.CheckField(validator.NotBlank(form.Content), "content", "This field can not be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expired", "This field must equal 1, 7, or 365")

	// use validator Valid() to validate if there are any errors
	if !form.Valid() {
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

	// use the Put() method to add a string value ("Snippet successfully created") and the key ("flash") to session
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// redirect to new snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
