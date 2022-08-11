package smtpcli

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/logger"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
)

type Config struct {
	Addr             string `yaml:"addr"`
	NoReplyEmailAddr string `yaml:"no_reply_email_addr"`
}

type Client struct {
	conf *Config
}

func NewClient(conf *Config) *Client {
	return &Client{
		conf: conf,
	}
}

func (c *Client) SendLetter(ctx context.Context, letter *service.Letter) error {
	logger.DefaultLogger().Info(c.conf.Addr)

	err := smtp.SendMail(
		c.conf.Addr,
		nil,
		c.conf.NoReplyEmailAddr,
		letter.EmailAddresses,
		[]byte(letter.Contents))
	if err != nil {
		return fmt.Errorf("cannot send email via smtp: %w", err)
	}

	return nil
}
