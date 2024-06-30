package queue_client

import (
	"context"
	"fmt"

	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/porebric/logger"
)

type Subscriber interface {
	Endpoint(ctx context.Context, request any) (response any, err error)
	DecodeRequest(context.Context, *Message) (request any, err error)
	Subject() string

	Decode(context.Context, []byte) (context.Context, *Message, error)
}

func Subscribe(ctx context.Context, conn Conn, ss ...Subscriber) error {
	natsConn, ok := conn.(*nats.Conn)
	if !ok {
		return fmt.Errorf("invalid nats connection structure")
	}

	for _, s := range ss {
		if _, err := natsConn.Subscribe(s.Subject(), makeSurveyHandler(natsConn, s)); err != nil {
			return fmt.Errorf("subscribe %s: %w", s.Subject(), err)
		}
		logger.Info(ctx, "subscribe", "subject", s.Subject())
	}

	return nil
}

func makeSurveyHandler(conn *nats.Conn, s Subscriber) nats.MsgHandler {
	opts := []kitnats.SubscriberOption{
		kitnats.SubscriberFinalizer(finalizer),
		kitnats.SubscriberErrorEncoder(encodeError),
	}

	return kitnats.NewSubscriber(
		s.Endpoint,
		func(ctx context.Context, msg *nats.Msg) (request any, err error) {
			ctx, req, err := s.Decode(ctx, msg.Data)
			if err != nil {
				return nil, fmt.Errorf("decode nats msg: %w", err)
			}
			return s.DecodeRequest(ctx, req)
		},
		encodeResponse,
		opts...,
	).ServeMsg(conn)
}
