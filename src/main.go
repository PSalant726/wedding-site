package main

import (
	"net/http"
	"text/template"
)

const (
	// ABOUT_PATH          = "/about"
	// ACCOMMODATIONS_PATH = "/accommodations"
	PREVIEW_PATH = "/preview"
	// RSVP_PATH           = "/rsvp"
)

var templates = template.Must(template.ParseFiles(getTemplateFiles()...))

func main() {
	http.HandleFunc("/", previewHandler)
	// http.HandleFunc(ABOUT_PATH, makeHandler(ABOUT_PATH))
	// http.HandleFunc(ACCOMMODATIONS_PATH, makeHandler(ACCOMMODATIONS_PATH))
	// http.HandleFunc(RSVP_PATH, makeHandler(RSVP_PATH))

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	startServer()
}

func previewHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate(w, PREVIEW_PATH); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// func makeHandler(path string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := renderTemplate(w, path); err != nil {
// 			http.Redirect(w, r, "/", http.StatusFound)
// 		}
// 	}
// }
