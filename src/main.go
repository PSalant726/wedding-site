package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Handler:      NewRouterWithRoutes(),
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  time.Minute,
	}

	log.Fatal(server.ListenAndServe())
}
