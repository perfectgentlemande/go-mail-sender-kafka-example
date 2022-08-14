package messagehandler

import (
	"context"
	"errors"
	"time"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
	"github.com/sirupsen/logrus"
)

type Config struct {
	RatePeriodMicroseconds int64 `yaml:"rate_period_microseconds"`
	RequestsPerPeriod      int64 `yaml:"requests_per_period"`
}

type MessageHandler interface {
	HandleMessages(ctx context.Context) error
}

func New(srvc *service.Service, log *logrus.Entry, conf *Config) MessageHandler {
	if conf.RatePeriodMicroseconds > 0 {
		return &messageHandlerWithRateLimiter{
			srvc:        srvc,
			log:         log,
			ticker:      time.NewTicker(time.Duration(conf.RatePeriodMicroseconds) * time.Microsecond),
			maxRequests: conf.RequestsPerPeriod,
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
			if errors.Is(err, service.ErrConnRefused) || errors.Is(err, service.ErrNoSuchHost) {
				return err
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
	srvc        *service.Service
	log         *logrus.Entry
	ticker      *time.Ticker
	maxRequests int64
}

func (mh *messageHandlerWithRateLimiter) HandleMessages(ctx context.Context) error {
	for {
		select {
		case <-mh.ticker.C:
			requestCounter := 0
			for requestCounter < int(mh.maxRequests) {
				m, err := mh.srvc.ReadLetter(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return nil
					}
					if errors.Is(err, service.ErrConnRefused) || errors.Is(err, service.ErrNoSuchHost) {
						return err
					}

					mh.log.WithError(err).Error("cannot read message")
					continue
				}

				err = mh.srvc.SendLetter(ctx, &m)
				if err != nil {
					mh.log.WithError(err).Error("cannot send message")
					requestCounter++
					continue
				}

				mh.log.WithFields(logrus.Fields{
					"emails":   m.EmailAddresses,
					"contents": m.Contents,
				}).Info("successfully sent message")
				requestCounter++
			}
		case <-ctx.Done():
			return nil
		}
	}
}
