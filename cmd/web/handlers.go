package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// use Header.Add() method to add custom header
	w.Header().Add("Server", "Go")

	// slice with the two template files
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// use template.ParseFiles() to read the files in the slice
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Execute method on the template set to write the template content
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
	}

}

// add snippetView handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// get wildcard, check if it's a positive integer
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// response
	_, err = fmt.Fprintf(w, "Display a specific snippet with ID: %d...", id)
	if err != nil {
		fmt.Println("Error:", err)
	}
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
	// user w.writeHeader to return a 201 using http constants
	w.WriteHeader(http.StatusCreated)
	// body as normal
	_, err := w.Write([]byte("Save a new snippet."))
	if err != nil {
		fmt.Println("Error:", err)
	}
}
