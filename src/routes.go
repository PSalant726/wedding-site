package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

const (
	PathAbout       = "/about"
	PathComm        = "/communicate"
	PathFAQ         = "/faq"
	PathHealth      = "/health"
	PathHome        = "/"
	PathPeople      = "/people"
	PathPreview     = "/preview"
	PathQuestion    = "/question"
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
		fs     = http.FileServer(http.Dir(RelativePathAssets))
		get    = router.Methods(http.MethodGet).Subrouter()
		getq   = get.Queries("address", "").Subrouter()
		post   = router.Methods(http.MethodPost).Subrouter()
	)

	router.NotFoundHandler = http.HandlerFunc(redirectHome)
	router.MethodNotAllowedHandler = http.HandlerFunc(redirectHome)
	router.HandleFunc(PathHealth, func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)
	router.PathPrefix("/" + RelativePathAssets).Handler(http.StripPrefix("/"+RelativePathAssets, fs))

	// GET requests
	get.Use(logRequest)
	get.HandleFunc(PathComm, makeHandler(PathComm))
	get.HandleFunc(PathHome, makeHandler(PathAbout))
	get.HandleFunc(PathPeople, makeHandler(PathPeople))
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
	post.HandleFunc(PathComm, commHandler)
	post.HandleFunc(PathPreview, subscribeHandler)
	post.HandleFunc(PathQuestion, questionHandler)
	post.HandleFunc(PathRSVP, rsvpHandler)
	post.HandleFunc(PathSubscribe, subscribeHandler)

	return router
}

func makeHandler(path string) http.HandlerFunc {
	p := NewPage(path)

	return func(w http.ResponseWriter, _ *http.Request) {
		if err := renderTemplate(w, path, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) error {
	templateFile := fmt.Sprintf("%s.html", tmpl)[1:]
	if err := templates.ExecuteTemplate(w, templateFile, p); err != nil {
		return fmt.Errorf("failed to execute template for file '%s': %w", templateFile, err)
	}

	return nil
}

func previewHandler(w http.ResponseWriter, _ *http.Request) {
	if err := templates.ExecuteTemplate(w, "preview.html", &Page{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func commHandler(w http.ResponseWriter, r *http.Request) {
	if pwd := r.FormValue("password"); pwd != os.Getenv("COMM_PASSWORD") {
		http.Error(w, "Incorrect Password", http.StatusUnauthorized)
		return
	}

	var (
		message          = r.FormValue("message")
		subscriberEmails = strings.Split(r.FormValue("emailAddresses"), ",")
		subscriberNames  = strings.Split(r.FormValue("names"), ",")
		subscriberList   = make(map[string]string)
	)

	for i, emailAddress := range subscriberEmails {
		subscriberList[emailAddress] = subscriberNames[i]
	}

	if err := emailSender.SendSubscriberCommunication(subscriberList, message); err != nil {
		http.Error(w, "Failed to send communication", http.StatusInternalServerError)
		log.Println(err)

		return
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

	if err := emailSender.SendSubscriberNotification(subscriber, true); err != nil {
		http.Error(w, "Failed to send subscriber notification", http.StatusInternalServerError)
		log.Println(err)

		if redirect {
			redirectHome(w, r)
		}

		return
	}

	msg := *NewSubscriberThankYouMessage(subscriber)
	if err := emailSender.SendHermesMessage(msg); err != nil {
		http.Error(w, "Failed to subscribe address: "+subscriber, http.StatusInternalServerError)
		log.Println(err)

		if redirect {
			redirectHome(w, r)
		}

		return
	}

	if redirect {
		redirectHome(w, r)
	}
}

func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query()["address"][0]

	if err := emailSender.SendSubscriberNotification(address, false); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
		log.Println(err)
		redirectHome(w, r)

		return
	}

	msg := *NewUnsubscribeConfirmationMessage(address)
	if err := emailSender.SendHermesMessage(msg); err != nil {
		http.Error(w, "Failed to unsubscribe address: "+address, http.StatusInternalServerError)
		log.Println(err)
		redirectHome(w, r)

		return
	}

	redirectHome(w, r)
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	var (
		senderName  = r.FormValue("name")
		senderEmail = r.FormValue("email")
		question    = r.FormValue("question")
	)

	if err := emailSender.SendQuestionNotification(senderName, senderEmail, question); err != nil {
		http.Error(w, "Failed to notify Phil & Rhiannon about your question. Please try again.", http.StatusInternalServerError)
		log.Println(err)

		return
	}

	msg := *NewQuestionReceivedMessage(senderName, senderEmail, question)
	if err := emailSender.SendHermesMessage(msg); err != nil {
		http.Error(w, "Failed to confirm receipt of your question. Please try again.", http.StatusInternalServerError)
		log.Println(err)
	}
}

func rsvpHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Endpoint not configured", http.StatusInternalServerError)
}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, PathHome, http.StatusPermanentRedirect)
}
