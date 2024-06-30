package queue_client

import (
	"context"
	"fmt"
)

type publisher[M Message] struct {
	conn Conn
}

func NewPublisher[M Message](conn Conn) Publisher[M] {
	return &publisher[M]{
		conn: conn,
	}
}

func (p *publisher[M]) Publish(_ context.Context, msg M) error {
	jsonData, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("nats marshal request: %w", err)
	}

	if err = p.conn.Publish(msg.GetSubject(), jsonData); err != nil {
		return fmt.Errorf("nats publish msg: %w", err)
	}

	return nil
}

func (p *publisher[M]) PublishWithResponse(ctx context.Context, msg M) error {
	jsonData, err := msg.MarshalBinary()
	if err != nil {
		return fmt.Errorf("nats marshal request: %w", err)
	}

	resp, err := p.conn.RequestWithContext(ctx, msg.GetSubject(), jsonData)
	if err != nil {
		return fmt.Errorf("nats publish msg: %w", err)
	}

	if err = msg.SetResponse(resp.Data); err != nil {
		return fmt.Errorf("nats set messages: %w", err)
	}

	return nil
}
