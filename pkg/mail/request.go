package mail

import (
	netmail "net/mail"
	"strings"
)

type Request struct {
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
	ReplyTo string `json:"replyTo"`
}

func (r *Request) Normalize() {
	r.Subject = strings.TrimSpace(r.Subject)
	r.Text = strings.TrimSpace(r.Text)
	r.HTML = strings.TrimSpace(r.HTML)
	r.ReplyTo = strings.TrimSpace(r.ReplyTo)
}

func (r Request) Valid() bool {
	if r.Subject == "" || strings.ContainsAny(r.Subject, "\r\n") {
		return false
	}

	if r.Text == "" && r.HTML == "" {
		return false
	}

	if r.ReplyTo == "" {
		return true
	}

	if strings.ContainsAny(r.ReplyTo, "\r\n") {
		return false
	}

	_, err := netmail.ParseAddress(r.ReplyTo)
	return err == nil
}
