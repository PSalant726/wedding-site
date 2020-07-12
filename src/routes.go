package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	PathAbout       = "/about"
	PathDetails     = "/details"
	PathFAQ         = "/faq"
	PathHealth      = "/health"
	PathHome        = "/"
	PathPreview     = "/preview"
	PathRegistry    = "/registry"
	PathRSVP        = "/rsvp"
	PathSchedule    = "/schedule"
	PathSubscribe   = "/subscribe"
	PathTravel      = "/travel"
	PathUnsubscribe = "/unsubscribe"

	RelativePathAssets = "assets/"
)

var (
	templates   = template.Must(template.ParseGlob("./assets/html/*"))
	emailSender = NewGmailUser("no-reply@rhiphilwedding.com", os.Getenv("GMAIL_PASSWORD"))
)

func NewRouterWithRoutes() *mux.Router {
	var (
		router = mux.NewRouter()
		get    = router.Methods(http.MethodGet).Subrouter()
		getq   = get.Queries("address", "").Subrouter()
		post   = router.Methods(http.MethodPost).Subrouter()
		fs     = http.FileServer(http.Dir(RelativePathAssets))
	)

	router.HandleFunc(PathHealth, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)
	router.PathPrefix("/" + RelativePathAssets).Handler(http.StripPrefix("/"+RelativePathAssets, fs))

	// GET requests
	get.Use(logRequest)
	get.HandleFunc(PathHome, makeHandler(PathAbout))
	get.HandleFunc(PathDetails, makeHandler(PathDetails))
	get.HandleFunc(PathFAQ, makeHandler(PathFAQ))
	get.HandleFunc(PathRegistry, makeHandler(PathRegistry))
	get.HandleFunc(PathRSVP, makeHandler(PathRSVP))
	get.HandleFunc(PathSchedule, makeHandler(PathSchedule))
	get.HandleFunc(PathTravel, makeHandler(PathTravel))
	get.HandleFunc(PathPreview, previewHandler)

	// GET requests with ?address=...
	getq.Use(logRequest)
	getq.HandleFunc(PathSubscribe, subscribeHandler)
	getq.HandleFunc(PathUnsubscribe, unsubscribeHandler)

	// POST requests
	post.Use(logRequest)
	post.HandleFunc(PathPreview, subscribeHandler)
	post.HandleFunc(PathSubscribe, subscribeHandler)

	return router
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

func previewHandler(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "preview.html", &Page{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		redirect   bool
		subscriber string
	)

	if r.Method == http.MethodGet {
		subscriber = r.URL.Query()["address"][0]
		redirect = true
	} else {
		subscriber = r.FormValue("email")
	}

	if err := emailSender.SendNotification(subscriber, true); err != nil {
		http.Error(w, "Failed to send subscriber notification", http.StatusInternalServerError)
		return
	}

	msg := *NewSubscriberThankYouMessage(subscriber)
	if err := emailSender.SendHermesMessage(msg); err != nil {
		http.Error(w, "Failed to subscribe address: "+subscriber, http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, PathHome, http.StatusPermanentRedirect)
	}
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query()["address"][0]

	if err := emailSender.SendNotification(address, false); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
	}

	msg := *NewUnsubscribeConfirmationMessage(address)
	if err := emailSender.SendHermesMessage(msg); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
	}

	http.Redirect(w, r, PathHome, http.StatusPermanentRedirect)
}
