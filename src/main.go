package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const (
	// ABOUT_PATH          = "/about"
	// ACCOMMODATIONS_PATH = "/accommodations"
	PREVIEW_PATH = "/preview"
	// RSVP_PATH           = "/rsvp"
)

var templates = template.Must(template.ParseGlob("./assets/html/*"))

func main() {
	r := mux.NewRouter()

	get := r.Methods("GET").Subrouter()
	get.HandleFunc("/", makeHandler(PREVIEW_PATH))
	get.HandleFunc(PREVIEW_PATH, makeHandler(PREVIEW_PATH))
	get.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	// get.HandleFunc(ABOUT_PATH, makeHandler(ABOUT_PATH))
	// get.HandleFunc(ACCOMMODATIONS_PATH, makeHandler(ACCOMMODATIONS_PATH))
	// get.HandleFunc(RSVP_PATH, makeHandler(RSVP_PATH))

	fs := http.FileServer(http.Dir("assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	server := &http.Server{
		Handler:      r,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
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
