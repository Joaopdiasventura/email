package mail

import "strings"

type Message struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
	ReplyTo string
}

func BuildMessage(msg Message) []byte {
	var builder strings.Builder

	builder.WriteString("From: ")
	builder.WriteString(msg.From)
	builder.WriteString("\r\n")
	builder.WriteString("To: ")
	builder.WriteString(msg.To)
	builder.WriteString("\r\n")
	builder.WriteString("Subject: ")
	builder.WriteString(msg.Subject)
	builder.WriteString("\r\n")

	if msg.ReplyTo != "" {
		builder.WriteString("Reply-To: ")
		builder.WriteString(msg.ReplyTo)
		builder.WriteString("\r\n")
	}

	if msg.HTML != "" {
		builder.WriteString("MIME-Version: 1.0\r\n")
		builder.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.HTML)
	} else {
		builder.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		builder.WriteString("\r\n")
		builder.WriteString(msg.Text)
	}

	return []byte(builder.String())
}
