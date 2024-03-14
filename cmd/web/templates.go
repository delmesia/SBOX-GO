package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/delmesia/snippet/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// humanDate function which will return a formatted
// string representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 3:04 PM")
}

// This will initialize a template.FuncMap object and store
// it in a local variable. This is a string-keyed map which will
// act as a lookup between names of custom template functions
// and the function themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// the function filepath.Glob() will return a slice of all the filepath
	// that matches the given pattern like "/Users/del/Dev/SBOX-GO/ui/html/pages/*.html"he
	// this will essentially give a slice of all the filepath for the application "page" templates
	// like: [/Users/del/Dev/SBOX-GO/ui/html/pages/view.html, /Users/del/Dev/SBOX-GO/ui/html/pages/view.html]
	pages, err := filepath.Glob("/Users/del/Dev/SBOX-GO/ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through the filepath one-by-one.
	for _, page := range pages {

		// Extract the file name (like "home.html") from the full filepath
		// and assign it to a new variable
		name := filepath.Base(page)
		// the template.FuncMap must be registered with the template set
		// before calling the ParseFiles() method. This means we have to use
		// template.New() to create an empty template set, use the Funcs() method
		// to register the template.FuncMap, and then parse the files as normal
		ts, err := template.New(name).Funcs(functions).ParseFiles("/Users/del/Dev/SBOX-GO/ui/html/base.html")
		if err != nil {
			return nil, err
		}
		/*// Parse the base template into a template set
		ts, err := template.ParseFiles("/Users/del/Dev/SBOX-GO/ui/html/base.html")
		if err != nil {
			return nil, err
		}*/
		// Call ParseGlob() *on this template*  to add any partials
		ts, err = ts.ParseGlob("/Users/del/Dev/SBOX-GO/ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Add the template set to the map, using the name of the page
		// (like "home.html") as the key
		cache[name] = ts
	}
	// return the map
	return cache, nil
}
