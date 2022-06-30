package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/logger"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/messagebroker"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/service"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/smtpcli"
	"golang.org/x/sync/errgroup"
)

func main() {
	log := logger.DefaultLogger()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	configPath := flag.String("c", "config.yaml", "path to your config")
	flag.Parse()

	conf, err := readConfig(*configPath)
	if err != nil {
		log.WithField("config_path", *configPath).WithError(err).Error("failed to read config")
		return
	}

	mBroker := messagebroker.New(conf.MessageBroker)
	smtpCli := smtpcli.NewClient(conf.SMTP)

	rungroup, ctx := errgroup.WithContext(ctx)

	srvc := service.New(mBroker, smtpCli)

	log.Info("starting server")
	rungroup.Go(func() error {
		for {
			m, err := srvc.ReadLetter(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}

				// use srvc.ReadAndSendLetters() to abort after the first read message error
				log.WithError(err).Error("cannot read message")
				continue
			}

			err = srvc.SendLetter(ctx, &m)
			if err != nil {
				// use srvc.ReadAndSendLetters() to abort after the first SMTP error
				log.WithError(err).Error("cannot send message")
			}
		}
	})

	rungroup.Go(func() error {
		<-ctx.Done()

		if err := mBroker.Close(); err != nil {
			return fmt.Errorf("failed to close messageBroker: %w", err)
		}

		return nil
	})

	err = rungroup.Wait()
	if err != nil {
		log.WithError(err).Error("run group exited because of error")
		return
	}

	log.Info("server exited properly")
}
