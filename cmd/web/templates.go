package main

import (
	"html/template"
	"path/filepath"

	"snippet.lciamp.xyz/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// use glob to get all templates from path
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through page filepaths
	for _, page := range pages {
		// extract name from file path
		name := filepath.Base(page)

		// create slice of filepaths
		files := []string{
			"ui/html/base.tmpl",
			"ui/html/partials/nav.tmpl",
			page,
		}

		// parse files into template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		// add template set to map
		cache[name] = ts
	}
	// return the map
	return cache, nil
}
