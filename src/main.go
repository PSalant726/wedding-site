package main

import (
	"log"
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

type Page struct {
	Title string
	Email []byte
}

func main() {
	http.HandleFunc("/", previewHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	// http.HandleFunc(ABOUT_PATH, makeHandler(ABOUT_PATH))
	// http.HandleFunc(ACCOMMODATIONS_PATH, makeHandler(ACCOMMODATIONS_PATH))
	// http.HandleFunc(RSVP_PATH, makeHandler(RSVP_PATH))

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(":5000", nil))
}

func previewHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Preview"}

	if err := renderTemplate(w, PREVIEW_PATH, p); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// func makeHandler(path string) http.HandlerFunc {
// 	p := &Page{Title: strings.Title(path[1:])}

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := renderTemplate(w, path, p); err != nil {
// 			http.Redirect(w, r, "/", http.StatusFound)
// 		}
// 	}
// }
