package main

import "strings"

type Page struct {
	Title  string
	Assets map[string]string
	Links  map[string]string
}

func NewPage(path string) *Page {
	page := path[1:]

	return &Page{
		Title: strings.Title(page),
		Assets: map[string]string{
			"css":    RelativePathAssets + "css/",
			"images": RelativePathAssets + "images/",
			"js":     RelativePathAssets + "js/",
		},
		Links: map[string]string{
			"faq":       PathFAQ,
			"home":      PathHome,
			"people":    PathPeople,
			"registry":  PathRegistry,
			"rsvp":      PathRSVP,
			"schedule":  PathSchedule,
			"subscribe": PathSubscribe,
			"travel":    PathTravel,
		},
	}
}
