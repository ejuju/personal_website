package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/ejuju/personal_website/app"
)

func main() {
	// Init and run HTTP server
	server := &http.Server{
		Handler: newHTTPHandler(),
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

//go:embed all:static
var staticFilesFS embed.FS

func newHTTPHandler() http.Handler {
	router := pat.New()

	// Serve pages
	router.Add(http.MethodGet, "/", app.ServeHomePage())
	router.Add(http.MethodGet, "/contact", app.ServeContactPage())
	router.Add(http.MethodGet, "/resume", app.ServeResumePage())
	router.Add(http.MethodGet, "/legal", app.ServeLegalPage())

	// Serve static files
	fsys, err := fs.Sub(staticFilesFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	router.NotFound = http.FileServer(http.FS(fsys))

	return router
}
