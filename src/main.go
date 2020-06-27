package main

import (
	"fmt"
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

var templates = template.Must(template.ParseGlob("./assets/html/*"))

func main() {
	http.HandleFunc("/", makeHandler(PREVIEW_PATH))
	// http.HandleFunc(ABOUT_PATH, makeHandler(ABOUT_PATH))
	// http.HandleFunc(ACCOMMODATIONS_PATH, makeHandler(ACCOMMODATIONS_PATH))
	// http.HandleFunc(RSVP_PATH, makeHandler(RSVP_PATH))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Fatal(http.ListenAndServe(":5000", nil))
}

func makeHandler(path string) http.HandlerFunc {
	p := NewPage(path)

	return func(w http.ResponseWriter, r *http.Request) {
		if err := renderTemplate(w, path, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) error {
	templateFile := fmt.Sprintf("%s.html", tmpl)[1:]
	if err := templates.ExecuteTemplate(w, templateFile, p); err != nil {
		return err
	}

	return nil
}
