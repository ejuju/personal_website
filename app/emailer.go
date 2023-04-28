package app

import (
	"fmt"
	"io"
)

type Email struct {
	Sender        string
	Recipient     string
	Subject       string
	PlainTextBody string
}

type Emailer interface {
	Send(*Email) error
}

type mockEmailer struct {
	w   io.Writer
	err error
}

func (e *mockEmailer) Send(email *Email) error {
	msg := fmt.Sprintf("New email: \n\tSender: %s\n\tRecipient: %s\n\tSubject: %s\n\tBody: %q\n",
		email.Sender,
		email.Recipient,
		email.Subject,
		email.PlainTextBody,
	)
	if e.err != nil {
		return e.err
	}
	_, err := e.w.Write([]byte(msg))
	return err
}
