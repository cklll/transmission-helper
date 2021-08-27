package main

import (
	"github.com/stretchr/testify/assert"
	"net/smtp"
	"testing"
)

func TestGetMailConfig(t *testing.T) {
	fakeAppConfig := getApplicationConfig("testdata/config/example.yaml")

	want := MailConfig{
		SmtpUser:        "test_user",
		SmtpPassword:    "test_pass",
		SmtpHost:        "test_host",
		SmtpPort:        "1025",
		IsNonSmtpSecure: true,
		SenderEmail:     "test_sender@email.com",
		SenderName:      "tranmission-helper",
	}
	result := MailNotifier{}.GetMailConfig(fakeAppConfig)

	assert.Equal(t, want, result)
}

func TestSendEncrypted(t *testing.T) {
	config := MailConfig{
		SmtpUser:        "test_user",
		SmtpPassword:    "test_password",
		SmtpHost:        "test_host",
		SmtpPort:        "1234",
		IsNonSmtpSecure: false,
		SenderEmail:     "test@sender.com",
		SenderName:      "test-sender",
	}
	subject := "test subject"
	message := "test message"
	recipients := []string{"test1@email.com", "test2@email.com"}

	mockSendMail := func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		wantMsg := []byte("From: test-sender <test@sender.com>\r\n" +
			"To: test1@email.com,test2@email.com\r\n" +
			"Subject: test subject\r\n" +
			"\r\n" +
			"test message")
		wantAuth := smtp.PlainAuth("", "test_user", "test_password", "test_host")

		assert.Equal(t, "test_host:1234", addr)
		assert.Equal(t, wantAuth, auth)
		assert.Equal(t, []string{"test1@email.com", "test2@email.com"}, to)
		assert.Equal(t, "test@sender.com", from)
		assert.Equal(t, wantMsg, msg)

		return nil
	}
	mailNotifier := MailNotifier{mockSendMail}
	mailNotifier.Send(config, subject, message, recipients)
}

func TestSendUnencrypted(t *testing.T) {
	config := MailConfig{
		SmtpUser:        "test_user",
		SmtpPassword:    "test_password",
		SmtpHost:        "test_host",
		SmtpPort:        "1234",
		IsNonSmtpSecure: true,
		SenderEmail:     "test@sender.com",
		SenderName:      "test-sender",
	}
	subject := "test subject"
	message := "test message"
	recipients := []string{"test1@email.com", "test2@email.com"}

	mockSendMail := func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
		wantAuth := UnencryptedAuth{smtp.PlainAuth("", "test_user", "test_password", "test_host")}
		assert.Equal(t, wantAuth, auth)

		return nil
	}
	mailNotifier := MailNotifier{mockSendMail}
	mailNotifier.Send(config, subject, message, recipients)
}
