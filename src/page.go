package main

import "strings"

type Page struct {
	Title string
	JS    string
}

func NewPage(path string) *Page {
	page := path[1:]

	return &Page{
		Title: strings.Title(page),
		JS:    page,
	}
}
