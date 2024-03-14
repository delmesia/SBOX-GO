package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// This will pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Since secureHeaders is just a function, and the function returns http.Handler
	// we don't need to do anything else.
	return secureHeaders(mux)
}
