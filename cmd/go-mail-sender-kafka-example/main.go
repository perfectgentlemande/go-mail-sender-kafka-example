package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer cancel()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		GroupID: "consumer-email",
		Topic:   "test",
	})

	rungroup, ctx := errgroup.WithContext(ctx)

	log.Println("starting server")
	rungroup.Go(func() error {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return nil
				}

				return fmt.Errorf("cannot read message: %w", err)
			}

			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	})

	rungroup.Go(func() error {
		<-ctx.Done()

		if err := r.Close(); err != nil {
			log.Fatal("failed to close reader:", err)
		}

		return nil
	})

	err := rungroup.Wait()
	if err != nil {
		log.Println("run group exited because of error", err)
		return
	}

	log.Println("server exited properly")
}
