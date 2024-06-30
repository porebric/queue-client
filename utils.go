package queue_client

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/porebric/logger"
)

type BaseResponse struct {
	Success  bool   `json:"success"`
	Response string `json:"messages"`
}

func encodeResponse(ctx context.Context, reply string, conn *nats.Conn, resp interface{}) error {
	b, err := json.Marshal(resp)
	if err != nil {
		logger.Error(ctx, err, "failed to marshal the messages data")
		return err
	}
	err = conn.Publish(reply, b)
	if err != nil {
		logger.Error(ctx, err, "failed to publish the reply")
		return err
	}
	return nil
}

func encodeError(ctx context.Context, err error, reply string, conn *nats.Conn) {
	b, err := json.Marshal(err)
	if err != nil {
		logger.Error(ctx, err, "failed to encode an error")
		return
	}
	_ = conn.Publish(reply, b)
}

func finalizer(ctx context.Context, msg *nats.Msg) {
	logger.Info(ctx, "request", "sub", msg.Subject, msg.Header.Get("sender"))
}
