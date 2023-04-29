package app

import (
	"bytes"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bmizerany/pat"
)

func NewHTTPHandler() http.Handler {
	// Init dependencies
	db := newInMemoryDB()
	emailer := &mockEmailer{w: os.Stdout, err: nil}
	router := pat.New()

	// Serve pages and forms
	router.Add(http.MethodGet, "/", prerenderAndServePage("home.gohtml", nil))
	router.Add(http.MethodGet, "/contact", prerenderAndServePage("contact.gohtml", nil))
	router.Add(http.MethodGet, "/contact_success", prerenderAndServePage("contact_success.gohtml", nil))
	router.Add(http.MethodPost, "/contact_form", handleContactForm(db, emailer))
	router.Add(http.MethodGet, "/resume", prerenderAndServePage("resume.gohtml", resumeTmplData))
	router.Add(http.MethodGet, "/legal", prerenderAndServePage("legal.gohtml", nil))

	// Serve static files
	fsys, err := fs.Sub(staticFilesFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	router.NotFound = http.FileServer(http.FS(fsys))

	// Wrap middleware
	out := newAnalyticsMiddleware(db)(router)
	out = newRecoveryMiddleware()(out)

	// Return HTTP handler
	return out
}

//go:embed all:ui
var uiFS embed.FS

//go:embed all:static
var staticFilesFS embed.FS

// Layout template files used for each page.
var layoutTmpls = []string{
	"ui/_layout.gohtml",
	"ui/_css.gohtml",
	"ui/_header.gohtml",
	"ui/_footer.gohtml",
}

func prerenderAndServePage(pageName string, data map[string]any) http.HandlerFunc {
	tmpl := template.Must(template.ParseFS(uiFS, append(layoutTmpls, "ui/"+pageName)...))
	buf := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(buf, "page_layout", data)
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) { w.Write(buf.Bytes()) }
}

var errPageTmpl = template.Must(template.ParseFS(uiFS, append(layoutTmpls, "ui/_error.gohtml")...))

func respondErrorPage(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	err := errPageTmpl.ExecuteTemplate(w, "page_layout", map[string]any{
		"Status":       strconv.Itoa(status),
		"StatusText":   http.StatusText(status),
		"ErrorMessage": message,
	})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func newID(length int) string {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(buf)
}
