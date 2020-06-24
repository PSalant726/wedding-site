package main

import (
	"fmt"
	"net/http"
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
		templateFiles = append(templateFiles, fmt.Sprintf("./assets/html%s.html", path))
	}

	return templateFiles
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) error {
	templateFile := fmt.Sprintf("%s.html", tmpl)[1:]
	if err := templates.ExecuteTemplate(w, templateFile, p); err != nil {
		return err
	}

	return nil
}
