package main

import "strings"

type Page struct {
	Title string
	Links map[string]string
}

func NewPage(path string) *Page {
	page := path[1:]

	return &Page{
		Title: strings.Title(page),
		Links: map[string]string{
			"details":   PathDetails,
			"faq":       PathFAQ,
			"home":      PathHome,
			"registry":  PathRegistry,
			"rsvp":      PathRSVP,
			"schedule":  PathSchedule,
			"subscribe": PathSubscribe,
			"travel":    PathTravel,
		},
	}
}
