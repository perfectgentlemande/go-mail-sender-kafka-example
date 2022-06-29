package service

import (
	"context"
	"errors"
	"fmt"
)

type Letter struct {
	EmailAddresses []string `json:"email_addresses"`
	Contents       string   `json:"contents"`
}

type Service struct {
	Broker  MessageBroker
	SMTPCli SMTPCli
}

func New(br MessageBroker, smtpCli SMTPCli) *Service {
	return &Service{
		Broker:  br,
		SMTPCli: smtpCli,
	}
}

func (s *Service) SendLetter(ctx context.Context, letter *Letter) error {
	return s.SMTPCli.SendLetter(ctx, letter)
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

		err = s.SendLetter(ctx, &m)
		if err != nil {
			return fmt.Errorf("cannot send message: %w", err)
		}
	}
}
