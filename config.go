package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ApplicationConfig struct {
	TransmissionRemote struct {
		Username string
		Password string
	} `yaml:"transmission_remote"`
	Smtp struct {
		Host        string
		Port        string
		NonSecure   string `yaml:"non_secure"`
		User        string
		Pass        string
		SenderEmail string `yaml:"sender_email"`
	}

	EmailRecipients []string `yaml:"email_recipients"`
}

func getApplicationConfig(configPath string) ApplicationConfig {
	config := ApplicationConfig{}

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err.Error())
		panic("Failed to load application config file.")
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err.Error())
		panic("Failed to unmarshal yaml file.")
	}

	return config
}
