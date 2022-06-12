package service

import "context"

type Letter struct {
	EmailAddress string
	Contents     string
}

type Service struct {
	Broker Broker
}

type LettersBroker interface {
	GetLetter(ctx context.Context) (Letter, error)
}

func (s *Service) ReadLetter(ctx context.Context) (Letter, error) {
	return s.Broker.ReadLetter(ctx)
}