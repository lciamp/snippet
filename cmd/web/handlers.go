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
	// REMOVED BECAUSE ADDED TO newTemplateData func
	//flash := app.sessionManager.PopString(r.Context(), "flash")

	// call newTemplateData helper
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// add the flash
	// data.Flash = flash

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

// user login handlers
type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// zero valued form
	var form userSignupForm

	// parse form data
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate using helper functions
	form.CheckField(validator.NotBlank(form.Name), "name", "This field is required")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRx), "email", "This must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters")

	// if any errors redisplay with a 422
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
	}

	// try to put new user in DB
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldErrors("email", "Email address already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// flash signup worked
	app.sessionManager.Put(r.Context(), "flash", "Signup Successful. Please login")

	// redirect to login
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
	//fmt.Fprintln(w, "Display login form")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validate the form. Check that email and pw are there and email matches the regex
	form.CheckField(validator.NotBlank(form.Email), "email", "This field is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRx), "password", "This should eb a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// check for valid credentials
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldErrors("Email or Password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// use renew token method on current session ID
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// add the ID of the current user to session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// redirect to create snippet page
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "logout user")
}
