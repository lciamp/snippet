package main

import (
	"fmt"
	"log"
	"net/http"
)

// home handler function
func home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from Snippet!"))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// add snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display a specific snippet..."))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// add a snippet handler function
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display form for creating a new snippet..."))
	if err != nil {
		fmt.Println("Error:", err)
	}

}

// main function
func main() {
	// create a servemux and register the home function
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	// register new handlers
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("starting on: 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
