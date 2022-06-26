package service

import (
	"context"
)

type MessageBroker interface {
	ReadLetter(ctx context.Context) (Letter, error)
}
