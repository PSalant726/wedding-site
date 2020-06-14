package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func getTemplateFiles() []string {
	var (
		templateFiles []string
		templatePaths = []string{
			// ABOUT_PATH,
			// ACCOMMODATIONS_PATH,
			PREVIEW_PATH,
			// RSVP_PATH,
		}
	)

	for _, path := range templatePaths {
		templateFiles = append(templateFiles, fmt.Sprintf("./templates%s.html", path))
	}

	return templateFiles
}

func renderTemplate(w http.ResponseWriter, tmpl string) error {
	templateFile := fmt.Sprintf("%s.html", tmpl)[1:]
	if err := templates.ExecuteTemplate(w, templateFile, nil); err != nil {
		return err
	}

	return nil
}

func startServer() {
	port, ok := os.LookupEnv("PORT") // for heroku
	if ok {
		port = fmt.Sprint(":", port)
	} else {
		port = ":8080"
	}

	log.Fatal(http.ListenAndServe(port, nil))
}
