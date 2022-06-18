package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/logger"
	"github.com/perfectgentlemande/go-mail-sender-kafka-example/internal/messagebroker"
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
		log.WithField("config_path", *configPath).WithError(err).Fatal("failed to read config")
	}

	mBroker := messagebroker.NewBroker(conf.MessageBroker)
	rungroup, ctx := errgroup.WithContext(ctx)

	log.Println("starting server")
	rungroup.Go(func() error {
		for {
			m, err := mBroker.ReadLetter(ctx)
			if err != nil {
				return fmt.Errorf("cannot read message: %w", err)
			}

			fmt.Printf("read message: %s", m.Contents)
		}
	})

	rungroup.Go(func() error {
		<-ctx.Done()

		if err := mBroker.Close(); err != nil {
			log.Fatal("failed to close messageBroker:", err)
		}

		return nil
	})

	err = rungroup.Wait()
	if err != nil {
		log.Println("run group exited because of error", err)
		return
	}

	log.Println("server exited properly")
}
