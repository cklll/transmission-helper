package main

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// https://stackoverflow.com/a/11066064/2691976
// Cannot configure a SMTP server with TLS for dev. This is the workaround
type UnencryptedAuth struct {
	smtp.Auth
}

func (a UnencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

type MailConfig struct {
	SmtpUser        string
	SmtpPassword    string
	SmtpHost        string
	SmtpPort        string
	IsNonSmtpSecure bool
	SenderEmail     string
	SenderName      string
}

type MailNotifier struct {
	SendMail func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

func (notifier MailNotifier) GetMailConfig() MailConfig {
	_, isNonSmtpSecure := os.LookupEnv("TH_SMTP_NON_SECURE")
	return MailConfig{
		SmtpUser:        os.Getenv("TH_SMTP_USER"),
		SmtpPassword:    os.Getenv("TH_SMTP_PASS"),
		SmtpHost:        os.Getenv("TH_SMTP_HOST"),
		SmtpPort:        os.Getenv("TH_SMTP_PORT"),
		IsNonSmtpSecure: isNonSmtpSecure,
		SenderEmail:     os.Getenv("TH_SMTP_SENDER_EMAIL"),
		SenderName:      "tranmission-helper",
	}
}

func (notifier MailNotifier) Send(config MailConfig, subject string, message string, recipients []string) error {
	plainAuth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPassword, config.SmtpHost)
	addr := fmt.Sprintf("%v:%v", config.SmtpHost, config.SmtpPort)
	var err error

	body := fmt.Sprintf(("From: %v <%v>\r\n" +
		"To: %v\r\n" +
		"Subject: %v\r\n" +
		"\r\n" +
		"%v"), config.SenderName, config.SenderEmail, strings.Join(recipients, ","), subject, message)

	if config.IsNonSmtpSecure {
		auth := UnencryptedAuth{plainAuth}
		err = notifier.SendMail(addr, auth, config.SenderEmail, recipients, []byte(body))
	} else {
		auth := plainAuth
		err = notifier.SendMail(addr, auth, config.SenderEmail, recipients, []byte(body))
	}

	return err
}
