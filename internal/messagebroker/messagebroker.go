package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	reader *kafka.Reader
}
type Config struct {
	BrokerAddr string `yaml:"broker_addr"`
	GroupID    string `yaml:"group_id"`
	Topic      string `yaml:"topic"`
	Key        string `yaml:"key"`
}

func New(conf *Config) *Broker {
	fmt.Println(conf.BrokerAddr)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{conf.BrokerAddr},
		GroupID: conf.GroupID,
		Topic:   conf.Topic,
	})

	return &Broker{
		reader: r,
	}
}

func (b *Broker) ReadLetter(ctx context.Context) (service.Letter, error) {
	m, err := b.reader.ReadMessage(ctx)
	if err != nil {
		return service.Letter{}, fmt.Errorf("cannot read message: %w", err)
	}

	lt := service.Letter{}
	err = json.Unmarshal(m.Value, &lt)
	if err != nil {
		return service.Letter{}, fmt.Errorf("cannot unmarshal message: %w", err)
	}

	return lt, nil
}

func (b *Broker) Close() error {
	return b.reader.Close()
}
