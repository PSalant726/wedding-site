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

const (
	// ABOUT_PATH          = "/about"
	// ACCOMMODATIONS_PATH = "/accommodations"
	PREVIEW_PATH = "/preview"
	// RSVP_PATH           = "/rsvp"
)

var (
	templates   = template.Must(template.ParseGlob("./assets/html/*"))
	emailSender = NewGmailUser("no-reply@rhiphilwedding.com", os.Getenv("GMAIL_PASSWORD"))
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)

	get := r.Methods(http.MethodGet).Subrouter()
	get.Use(logRequest)
	get.HandleFunc("/", makeHandler(PREVIEW_PATH))
	get.HandleFunc(PREVIEW_PATH, makeHandler(PREVIEW_PATH))
	// get.HandleFunc(ABOUT_PATH, makeHandler(ABOUT_PATH))
	// get.HandleFunc(ACCOMMODATIONS_PATH, makeHandler(ACCOMMODATIONS_PATH))
	// get.HandleFunc(RSVP_PATH, makeHandler(RSVP_PATH))
	get.HandleFunc("/subscribe", subscribeHandler).Queries("address", "")
	get.HandleFunc("/unsubscribe", unsubscribeHandler).Queries("address", "")

	post := r.Methods(http.MethodPost).Subrouter()
	post.Use(logRequest)
	post.HandleFunc(PREVIEW_PATH, subscribeHandler)

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

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		redirect   bool
		subscriber = r.FormValue("email")
	)

	if subscriber == "" {
		subscriber = r.URL.Query()["address"][0]
		redirect = true
	}

	if err := emailSender.SendNotification(subscriber, true); err != nil {
		http.Error(w, "Failed to send subscriber notification", http.StatusInternalServerError)
		return
	}

	message, subject := NewSubscriberThankYouMessage(subscriber)
	if err := emailSender.SendHermesMessage(subscriber, subject, message); err != nil {
		http.Error(w, "Failed to subscribe address: "+subscriber, http.StatusInternalServerError)
	}

	if redirect {
		http.Redirect(w, r, PREVIEW_PATH, http.StatusPermanentRedirect)
	}
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query()["address"][0]

	if err := emailSender.SendNotification(address, false); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
	}

	message, subject := NewUnsubscribeConfirmationMessage(address)
	if err := emailSender.SendHermesMessage(address, subject, message); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
