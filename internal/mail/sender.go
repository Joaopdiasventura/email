package mail

import (
	"crypto/tls"
	"net/smtp"
	"strconv"

	"github.com/Joaopdiasventura/email/internal/config"
)

type Sender struct {
	Config config.SMTPConfig
}

func NewSender(cfg config.SMTPConfig) Sender {
	return Sender{
		Config: cfg,
	}
}

func (s Sender) Send(to string, req Request) error {
	message := BuildMessage(Message{
		From:    s.Config.From,
		To:      to,
		Subject: req.Subject,
		Text:    req.Text,
		HTML:    req.HTML,
		ReplyTo: req.ReplyTo,
	})

	address := s.Config.Host + ":" + strconv.Itoa(s.Config.Port)
	auth := smtp.PlainAuth("", s.Config.User, s.Config.Pass, s.Config.Host)

	if s.Config.Secure {
		return s.sendSecure(address, auth, s.Config.From, to, message)
	}

	return smtp.SendMail(address, auth, s.Config.From, []string{to}, message)
}

func (s Sender) sendSecure(address string, auth smtp.Auth, from string, to string, message []byte) error {
	conn, err := tls.Dial("tcp", address, &tls.Config{
		ServerName: s.Config.Host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.Config.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	if err := client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}
