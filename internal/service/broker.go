package service

import (
	"context"
)

type Broker interface {
	ReadLetter(ctx context.Context) (Letter, error)
}
