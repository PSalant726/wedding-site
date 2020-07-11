package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

var (
	templates   = template.Must(template.ParseGlob("./assets/html/*"))
	emailSender = NewGmailUser("no-reply@rhiphilwedding.com", os.Getenv("GMAIL_PASSWORD"))
)

func main() {
	var (
		r      = mux.NewRouter()
		server = &http.Server{
			Handler:      r,
			Addr:         ":5000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
	)

	AddRoutes(r)
	log.Fatal(server.ListenAndServe())
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + ": " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
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
