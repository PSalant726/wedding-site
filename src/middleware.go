package main

import (
	"log"
	"net/http"
)

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
