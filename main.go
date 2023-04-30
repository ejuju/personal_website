package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ejuju/personal_website/app"
)

func main() {
	// Init and run HTTP server
	server := &http.Server{
		Handler:        app.NewHTTPHandler(os.Getenv("MODE") != "PROD"),
		Addr:           ":8080",
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		IdleTimeout:    time.Second,
		MaxHeaderBytes: 8000,
	}
	log.Println("starting HTTP server at address", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
