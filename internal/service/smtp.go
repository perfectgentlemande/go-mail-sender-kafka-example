package service

import "context"

type SMTPCli interface {
	SendLetter(ctx context.Context, letter *Letter) error
}
