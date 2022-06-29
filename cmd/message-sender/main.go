package main

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/logger"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
	"github.com/segmentio/kafka-go"
)

// this is the example message sender
func main() {
	log := logger.DefaultLogger()

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "mailsender", 0)
	if err != nil {
		log.WithError(err).Fatal("cannot create connection")
	}
	defer conn.Close()

	messages := make([]kafka.Message, 5)
	for i := range messages {
		uuidVal := uuid.NewString()

		msgPayload := service.Letter{
			EmailAddresses: []string{uuidVal + "@test.com"},
			Contents:       "Hello, " + uuidVal,
		}

		data, err := json.Marshal(msgPayload)
		if err != nil {
			log.WithError(err).Error("cannot marshal message")
			return
		}

		messages[i] = kafka.Message{
			Key:   []byte("message"),
			Value: data,
		}
	}

	_, err = conn.WriteMessages(messages...)
	if err != nil {
		log.WithError(err).Error("cannot write messages")
	}
}
