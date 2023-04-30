package main

import (
	"net/http"
	"time"

	"github.com/ejuju/personal_website/app"
)

func main() {
	// Init and run HTTP server
	server := &http.Server{
		Handler:        app.NewHTTPHandler(true),
		Addr:           ":8080",
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		IdleTimeout:    time.Second,
		MaxHeaderBytes: 8000,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
