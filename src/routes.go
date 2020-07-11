package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PathAbout       = "/about"
	PathHealth      = "/health"
	PathPreview     = "/preview"
	PathSubscribe   = "/subscribe"
	PathUnsubscribe = "/unsubscribe"

	RelativePathAssets = "assets/"
)

func AddRoutes(r *mux.Router) {
	var (
		get  = r.Methods(http.MethodGet).Subrouter()
		getq = get.Queries("address", "").Subrouter()
		post = r.Methods(http.MethodPost).Subrouter()
		fs   = http.FileServer(http.Dir(RelativePathAssets))
	)

	r.HandleFunc(PathHealth, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)
	r.PathPrefix("/" + RelativePathAssets).Handler(http.StripPrefix("/"+RelativePathAssets, fs))

	// GET requests
	get.Use(logRequest)
	get.HandleFunc("/", makeHandler(PathAbout))
	get.HandleFunc(PathPreview, previewHandler)

	// GET requests with ?address=...
	getq.Use(logRequest)
	getq.HandleFunc(PathSubscribe, subscribeHandler)
	getq.HandleFunc(PathUnsubscribe, unsubscribeHandler)

	// POST requests
	post.Use(logRequest)
	post.HandleFunc(PathPreview, subscribeHandler)
	post.HandleFunc(PathSubscribe, subscribeHandler)
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
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
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

	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}
