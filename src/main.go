package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	templateFiles = []string{
		// "./templates/about.html",
		// "./templates/accommodations.html",
		"./templates/preview.html",
		// "./templates/rsvp.html",
	}
	templates = template.Must(template.ParseFiles(templateFiles...))
)

func main() {
	http.HandleFunc("/", previewHandler)
	// http.HandleFunc("/", aboutHandler)
	// http.HandleFunc("/accommodations", accommodationsHandler)
	// http.HandleFunc("/rsvp", rsvpHandler)

	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	port, ok := os.LookupEnv("PORT") // for heroku
	if ok {
		port = fmt.Sprint(":", port)
		log.Fatal(http.ListenAndServe(port, nil))
		return
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func previewHandler(w http.ResponseWriter, r *http.Request) {
	if err := renderTemplate(w, "preview"); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := renderTemplate(w, "about"); err != nil {
// 		redirectHome(w, r)
// 	}
// }

// func accommodationsHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := renderTemplate(w, "accommodations"); err != nil {
// 		redirectHome(w, r)
// 	}
// }

// func rsvpHandler(w http.ResponseWriter, r *http.Request) {
// 	if err := renderTemplate(w, "rsvp"); err != nil {
// 		redirectHome(w, r)
// 	}
// }

func renderTemplate(w http.ResponseWriter, tmpl string) error {
	templateFile := fmt.Sprintf("%s.html", tmpl)
	if err := templates.ExecuteTemplate(w, templateFile, nil); err != nil {
		return err
	}

	return nil
}

// func redirectHome(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/", http.StatusFound)
// }
