package main

import (
	"fmt"
	"os"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/messagebroker"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/smtpcli"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MessageBroker  *messagebroker.Config  `yaml:"messagebroker"`
	SMTP           *smtpcli.Config        `yaml:"smtp"`
	MessageHandler *messagehandler.Config `yaml:"messagehandler"`
}

func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return config, nil
}
