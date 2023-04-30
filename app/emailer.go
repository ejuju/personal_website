package app

import (
	"fmt"
	"io"
	"net/smtp"
	"strconv"
	"strings"
)

type Email struct {
	From          string
	To            []string
	Subject       string
	PlainTextBody string
}

type Emailer func(*Email) error

func sendEmailToAdmin(config *Config, emailer Emailer, subject, msg string) error {
	return emailer(&Email{
		From:          config.SMTPSender,
		To:            []string{config.AdminEmailAddr},
		Subject:       subject,
		PlainTextBody: msg,
	})
}

func newMockEmailer(w io.Writer, err error) Emailer {
	return func(email *Email) error {
		msg := fmt.Sprintf("New email: \n\tFrom: %s\n\tTo: %s\n\tSubject: %s\n\tBody:\n\n%s\n",
			email.From,
			email.To,
			email.Subject,
			email.PlainTextBody,
		)
		if err != nil {
			return err
		}
		_, err := w.Write([]byte(msg))
		return err
	}
}

func newSMTPEmailer(config *Config) Emailer {
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	return func(email *Email) error {
		addr := config.SMTPHost + ":" + strconv.Itoa(config.SMTPPort)
		return smtp.SendMail(addr, auth, config.SMTPUsername, email.To, []byte(emailMessageStr(email)))
	}
}

// generates the message string that will be sent to the SMTP server
func emailMessageStr(e *Email) string {
	headerMap := map[string]string{
		"From":         e.From,
		"To":           strings.Join(e.To, "; "),
		"Subject":      e.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain",
	}
	header := ""
	for key, val := range headerMap {
		header += key + ":" + val + "\r\n"
	}
	body := e.PlainTextBody
	return header + "\r\n" + body + "\r\n"
}
