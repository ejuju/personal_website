package app

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"
)

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

func handleContactForm(config *Config, db DB, emailer Emailer) http.HandlerFunc {
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
			ID:           newID(32),
			CreatedAt:    time.Now(),
			EmailAddress: emailAddress.Address,
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
		err = sendEmailToAdmin(config, emailer, "New contact form submission", contactFormSubmission.String())
		if err != nil {
			log.Println(err)
			respondErrorPage(w, http.StatusInternalServerError, "failed to send notification email")
			return
		}

		// Send confirmation email to user
		err = emailer(&Email{
			Sender:        config.SMTPSender,
			Recipients:    []string{emailAddress.Address},
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