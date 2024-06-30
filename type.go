package nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

type Publisher[M Message] interface {
	PublishWithResponse(ctx context.Context, msg M) error
	Publish(_ context.Context, msg M) error
}

type Client interface {
	Conn() Conn
}

type Message interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error

	SetResponse(data []byte) error
	GetSubject() string
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.0 --name=Conn
type Conn interface {
	Publish(subj string, data []byte) error
	RequestWithContext(ctx context.Context, subj string, data []byte) (*nats.Msg, error)
	Close()
}
