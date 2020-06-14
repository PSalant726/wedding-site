package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
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
	t, _ := template.ParseFiles("./templates/preview.html")

	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./templates/about.html")
// 	t.Execute(w, nil)
// }

// func accommodationsHandler(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./templates/accommodations.html")
// 	t.Execute(w, nil)
// }

// func rsvpHandler(w http.ResponseWriter, r *http.Request) {
// 	t, _ := template.ParseFiles("./templates/rsvp.html")
// 	t.Execute(w, nil)
// }
