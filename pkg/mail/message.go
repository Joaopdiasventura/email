package mail

import (
	"errors"
	"fmt"
	"mime"
	netmail "net/mail"
	"strings"
)

type Message struct {
	From    string
	To      string
	Subject string
	Text    string
	HTML    string
	ReplyTo string
}

func BuildMessage(msg Message) ([]byte, error) {
	from, err := formatAddressHeader("From", msg.From)
	if err != nil {
		return nil, err
	}

	to, err := formatAddressHeader("To", msg.To)
	if err != nil {
		return nil, err
	}

	subject, err := formatHeaderValue("Subject", msg.Subject)
	if err != nil {
		return nil, err
	}

	var builder strings.Builder

	builder.WriteString("From: ")
	builder.WriteString(from)
	builder.WriteString("\r\n")
	builder.WriteString("To: ")
	builder.WriteString(to)
	builder.WriteString("\r\n")
	builder.WriteString("Subject: ")
	builder.WriteString(mime.QEncoding.Encode("UTF-8", subject))
	builder.WriteString("\r\n")

	if msg.ReplyTo != "" {
		replyTo, err := formatAddressHeader("Reply-To", msg.ReplyTo)
		if err != nil {
			return nil, err
		}

		builder.WriteString("Reply-To: ")
		builder.WriteString(replyTo)
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

	return []byte(builder.String()), nil
}

func formatAddressHeader(name string, value string) (string, error) {
	value, err := formatHeaderValue(name, value)
	if err != nil {
		return "", err
	}

	address, err := netmail.ParseAddress(value)
	if err != nil {
		return "", fmt.Errorf("invalid %s header: %w", name, err)
	}

	return address.String(), nil
}

func formatHeaderValue(name string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("empty %s header", name)
	}

	if strings.ContainsAny(value, "\r\n") {
		return "", fmt.Errorf("invalid %s header", name)
	}

	return value, nil
}

func envelopeAddress(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("empty envelope address")
	}

	address, err := netmail.ParseAddress(value)
	if err != nil {
		return "", fmt.Errorf("invalid envelope address: %w", err)
	}

	return address.Address, nil
}
