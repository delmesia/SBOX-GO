package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// TODO: unused parameter. will be using soon
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	// This will retrieve the appropriate template set from cache based on the page like
	// "home.html". If no entry exists in the cache with the provided name,
	// then create a new error and call the serverError() helper then return.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	// Initialize a new buffer.
	buf := new(bytes.Buffer)
	// Write the template to the buffer, instead of sending straight to
	// the http.ResponseWriter. If an error occured, call the serverError() helper
	// and return again.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// If the template is written to the buffer without any errors,
	// its safe to write the HTTP status code to the http.ResponseWriter
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter
	buf.WriteTo(w)
}
