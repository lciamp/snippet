package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippet.lciamp.xyz/internal/models"
	"snippet.lciamp.xyz/ui"
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

	// use glob to get all templates from embedded filesystem
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through page file paths
	for _, page := range pages {
		// extract name from file path
		name := filepath.Base(page)

		// slice of filepath patterns with templates we want to parse
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// use ParseF() instead of ParseFiles to parse the template files
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// add template set to map
		cache[name] = ts
	}
	// return the map
	return cache, nil
}
