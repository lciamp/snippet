package main

import (
	"html/template"
	"path/filepath"
	"snippet.lciamp.xyz/internal/models"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// human date function
func humanDate(t time.Time) string {
	return t.Format("Jan 02 2006 at 15:04")
}

// initialize funcMap
var functions = template.FuncMap{
	"humanDate": humanDate,
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
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
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
