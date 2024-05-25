package main

import (
	"log"
	"net/http"
)

// home handler function
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippet!"))
}

// main function
func main() {
	// create a servemux and register the home function
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("starting on: 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
