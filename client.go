package queue_client

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/porebric/logger"
)

type client struct {
	conn Conn
}

func NewNatsClient(ctx context.Context, host string, port int) (Client, error) {
	conn, err := nats.Connect(fmt.Sprintf("nats://%s:%d", host, port),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info(ctx, "reconnected to NATS server", "url", nc.ConnectedUrl())
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error(ctx, err, "disconnected from NATS server")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Error(ctx, nc.LastError(), "connection closed")
		}),
		nats.MaxReconnects(-1),
	)
	if err != nil {
		return nil, err
	}

	return &client{
		conn: conn,
	}, nil
}

func (h *client) Conn() Conn {
	return h.conn
}
