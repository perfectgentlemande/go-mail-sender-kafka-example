package messagehandler

import (
	"context"
	"errors"
	"time"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
	"github.com/sirupsen/logrus"
)

type Config struct {
	RateMicroseconds int64 `yaml:"rate_microseconds"`
}

type MessageHandler interface {
	HandleMessages(ctx context.Context) error
}

func New(srvc *service.Service, log *logrus.Entry, conf *Config) MessageHandler {
	if conf.RateMicroseconds > 0 {
		return &messageHandlerWithRateLimiter{
			srvc:   srvc,
			log:    log,
			ticker: time.NewTicker(time.Duration(conf.RateMicroseconds) * time.Microsecond),
		}
	}

	return &messageHandler{
		srvc: srvc,
		log:  log,
	}
}

type messageHandler struct {
	srvc *service.Service
	log  *logrus.Entry
}

func (mh *messageHandler) HandleMessages(ctx context.Context) error {
	for {
		m, err := mh.srvc.ReadLetter(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			// use srvc.ReadAndSendLetters() to abort after the first read message error
			mh.log.WithError(err).Error("cannot read message")
			continue
		}

		err = mh.srvc.SendLetter(ctx, &m)
		if err != nil {
			// use srvc.ReadAndSendLetters() to abort after the first SMTP error
			mh.log.WithError(err).Error("cannot send message")
		}

		mh.log.WithFields(logrus.Fields{
			"emails":   m.EmailAddresses,
			"contents": m.Contents,
		}).Info("successfully sent message")
	}
}

type messageHandlerWithRateLimiter struct {
	srvc   *service.Service
	log    *logrus.Entry
	ticker *time.Ticker
}

func (mh *messageHandlerWithRateLimiter) HandleMessages(ctx context.Context) error {
	for {
		select {
		case <-mh.ticker.C:
			m, err := mh.srvc.ReadLetter(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}

				mh.log.WithError(err).Error("cannot read message")
				continue
			}

			err = mh.srvc.SendLetter(ctx, &m)
			if err != nil {
				mh.log.WithError(err).Error("cannot send message")
			}

			mh.log.WithFields(logrus.Fields{
				"emails":   m.EmailAddresses,
				"contents": m.Contents,
			}).Info("successfully sent message")
		case <-ctx.Done():
			return nil
		}
	}
}
