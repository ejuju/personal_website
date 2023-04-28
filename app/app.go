package app

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/bmizerany/pat"
)

func NewHTTPHandler() http.Handler {
	// Init dependencies
	db := &inMemoryDB{contactFormSubmissions: map[string]*ContactFormSubmission{}}
	emailer := &mockEmailer{w: os.Stdout, err: errors.New("fake send error")}
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

	// Return HTTP handler
	return router
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

type ContactFormSubmission struct {
	ID           string
	CreatedAt    time.Time
	EmailAddress string
	Message      string
}

func (s *ContactFormSubmission) String() string {
	out := "ID: " + s.ID + "\n"
	out += "Created at: " + s.CreatedAt.Format(time.RFC3339) + "\n"
	out += "Email address: " + s.EmailAddress + "\n"
	out += "Message: " + s.Message + "\n"
	return out
}

func handleContactForm(db DB, emailer Emailer) http.HandlerFunc {
	maxMessageLength := 8000

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and validate form
		err := r.ParseForm()
		if err != nil {
			respondErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		emailAddress, err := mail.ParseAddress(r.FormValue("email_address"))
		if err != nil {
			respondErrorPage(w, http.StatusBadRequest, err.Error())
			return
		}
		message := r.FormValue("message")
		if len(message) > maxMessageLength {
			errmsg := fmt.Sprintf("message length (%d) is too long (max %d characters)", len(message), maxMessageLength)
			respondErrorPage(w, http.StatusBadRequest, errmsg)
			return
		}

		// Set id and created at timestamp
		contactFormSubmission := &ContactFormSubmission{
			ID:           time.Now().Format(time.RFC3339) + "/" + emailAddress.Address,
			CreatedAt:    time.Now(),
			EmailAddress: emailAddress.String(),
			Message:      message,
		}

		// Store message in DB
		err = db.NewContactFormSubmission(contactFormSubmission)
		if err != nil {
			log.Println(err)
			respondErrorPage(w, http.StatusInternalServerError, "failed to save to database")
			return
		}

		// Send notification email to admin
		err = emailer.Send(&Email{
			Sender:        "bot@juliensellier.com",
			Recipient:     "admin@juliensellier.com",
			Subject:       "New contact form submission",
			PlainTextBody: contactFormSubmission.String(),
		})
		if err != nil {
			log.Println(err)
			respondErrorPage(w, http.StatusInternalServerError, "failed to send notification email")
			return
		}

		// Send confirmation email to user
		err = emailer.Send(&Email{
			Sender:        "bot@juliensellier.com",
			Recipient:     contactFormSubmission.EmailAddress,
			Subject:       "Thank you for your message!",
			PlainTextBody: contactFormSubmission.String(),
		})
		if err != nil {
			log.Println(err)
			respondErrorPage(w, http.StatusInternalServerError, "failed to send confirmation email")
			return
		}

		// Redirect to success page
		http.Redirect(w, r, "/contact_success", http.StatusSeeOther)
	}
}
