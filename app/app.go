package app

import (
	"bytes"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"

	"github.com/bmizerany/pat"
)

//go:embed all:ui
var uiFS embed.FS

//go:embed all:static
var staticFilesFS embed.FS

// Layout template files used for each page.
var layoutTmpls = []string{
	"ui/_layout.gohtml",
	"ui/_header.gohtml",
	"ui/_footer.gohtml",
}

func NewHTTPHandler() http.Handler {
	router := pat.New()

	// Serve pages
	router.Add(http.MethodGet, "/", serveHomePage())
	router.Add(http.MethodGet, "/contact", serveContactPage())
	router.Add(http.MethodGet, "/resume", serveResumePage())
	router.Add(http.MethodGet, "/legal", serveLegalPage())

	// Serve static files
	fsys, err := fs.Sub(staticFilesFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	router.NotFound = http.FileServer(http.FS(fsys))

	return router
}

func prerenderPage(pageName string, data map[string]any) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(uiFS, append(layoutTmpls, "ui/"+pageName)...))
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, "page_layout", data)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) { w.Write(buf.Bytes()) }
}

func serveHomePage() http.HandlerFunc    { return prerenderPage("home.gohtml", nil) }
func serveContactPage() http.HandlerFunc { return prerenderPage("contact.gohtml", nil) }
func serveResumePage() http.HandlerFunc  { return prerenderPage("resume.gohtml", resumeTmplData) }
func serveLegalPage() http.HandlerFunc   { return prerenderPage("legal.gohtml", nil) }
