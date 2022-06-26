package service

import (
	"context"
	"errors"
	"fmt"
)

type Letter struct {
	EmailAddress string `json:"email_address"`
	Contents     string `json:"contents"`
}

type Service struct {
	Broker MessageBroker
}

func New(br MessageBroker) *Service {
	return &Service{
		Broker: br,
	}
}

func (s *Service) ReadLetter(ctx context.Context) (Letter, error) {
	return s.Broker.ReadLetter(ctx)
}

func (s *Service) ReadLetters(ctx context.Context) error {
	for {
		m, err := s.ReadLetter(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return fmt.Errorf("cannot read message: %w", err)
		}

		fmt.Println(m.Contents)
	}
}
