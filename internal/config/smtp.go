package config

import (
	"errors"
	"os"
	"strconv"
)

type SMTPConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	From   string
	Secure bool
}

func LoadSMTPConfig() (SMTPConfig, error) {
	portValue := os.Getenv("SMTP_PORT")

	if portValue == "" {
		portValue = "587"
	}

	port, err := strconv.Atoi(portValue)
	if err != nil {
		return SMTPConfig{}, errors.New("INVALID_SMTP_PORT")
	}

	cfg := SMTPConfig{
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		From: os.Getenv("SMTP_FROM"),
	}

	secureValue := os.Getenv("SMTP_SECURE")

	cfg.Secure = secureValue == "true" || port == 465

	if cfg.Host == "" || cfg.User == "" || cfg.Pass == "" || cfg.From == "" {
		return SMTPConfig{}, errors.New("MISSING_SMTP_ENV")
	}

	return cfg, nil
}
