package main

import (
	"log"
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestGetApplicationConfig(t *testing.T) {
	// not sure if this is a good idea.
	// if we want to turn off log, we probably want to do it using a global setup/teardown
	log.SetOutput(ioutil.Discard)

	t.Run("when everything is correct", func(t *testing.T) {
		gotConfig := getApplicationConfig("./testdata/config/example.yaml")
		wantNotifyEmails := []string{"test_receiver1@example.com", "test_receiver2@example.com"}

		assert.Equal(t, "test_tr_user", gotConfig.TransmissionRemote.Username)
		assert.Equal(t, "test_tr_password", gotConfig.TransmissionRemote.Password)

		assert.Equal(t, "test_host", gotConfig.Smtp.Host)
		assert.Equal(t, "1025", gotConfig.Smtp.Port)
		assert.Equal(t, wantNotifyEmails, gotConfig.EmailRecipients)
	})

	t.Run("when file path doesn't exists", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Reading non-existing file didn't panic")
			}
		}()

		getApplicationConfig("./404.yaml")
	})

	t.Run("when file isn't in yaml format", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Reading invalid yaml file didn't panic")
			}
		}()

		getApplicationConfig("./testdata/config/bad_format.yaml")
	})
}
