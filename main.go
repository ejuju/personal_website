package main

import (
	"net/http"

	"github.com/ejuju/personal_website/app"
)

func main() {
	// Init and run HTTP server
	server := &http.Server{
		Handler: app.NewHTTPHandler(),
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
