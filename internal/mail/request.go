package mail

import (
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
	return r.Subject != "" && (r.Text != "" || r.HTML != "")
}
