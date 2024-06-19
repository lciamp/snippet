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

		// parse base template file into a templateset
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse Glob to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// add page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to map
		cache[name] = ts
	}
	// return the map
	return cache, nil
}
