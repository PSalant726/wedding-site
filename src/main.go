package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Handler:      NewRouterWithRoutes(),
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string
		if _, ok := r.Header["X-Forwarded-For"]; ok {
			ip = r.Header["X-Forwarded-For"][0]
		} else {
			ip = r.RemoteAddr
		}

		log.Printf("(%s) %s: %s", ip, r.Method, r.URL)
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
	var redirect bool

	subscriber := r.FormValue("email")
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
