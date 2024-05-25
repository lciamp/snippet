package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// home handler function
func home(w http.ResponseWriter, r *http.Request) {
	// use Hdeader.Add() method to add custom header
	w.Header().Add("Server", "Go")
	_, err := w.Write([]byte("Hello from Snippet!"))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// add snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
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
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display form for creating a new snippet..."))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// add a snippet handler to POST snippet
func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// user w.writeHeader to return a 201 using http constants
	w.WriteHeader(http.StatusCreated)
	// body as normal
	_, err := w.Write([]byte("Save a new snippet."))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// main function
func main() {
	// create a servemux and register the home function
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // restrict for / only
	// register new handlers
	mux.HandleFunc("GET /snippet/view/{id}/{$}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting on: localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
