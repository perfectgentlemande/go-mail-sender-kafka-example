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
	Addr    string `yaml:"addr"`
	GroupID string `yaml:"group_id"`
	Topic   string `yaml:"topic"`
	Key     string `yaml:"key"`
}

func NewBroker(conf *Config) *Broker {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{conf.Addr},
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
		if err == context.Canceled {
			return service.Letter{}, nil
		}

		return service.Letter{}, fmt.Errorf("cannot read message: %w", err)
	}

	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

	lt := service.Letter{}
	err = json.Unmarshal(m.Value, &lt)
	if err != nil {
		return service.Letter{}, fmt.Errorf("cannot unmarshal message: %w", err)
	}

	return lt, nil
}