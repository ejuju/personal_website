package app

import (
	"bytes"
	"embed"
	"net/http"
	"text/template"
)

//go:embed all:ui
var uiFS embed.FS

var layoutTmpls = []string{
	"ui/_layout.gohtml",
	"ui/_header.gohtml",
	"ui/_footer.gohtml",
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

func ServeHomePage() http.HandlerFunc    { return prerenderPage("home.gohtml", nil) }
func ServeContactPage() http.HandlerFunc { return prerenderPage("contact.gohtml", nil) }
func ServeCVPage() http.HandlerFunc      { return prerenderPage("cv.gohtml", cvPageData) }
