package logit

import (
	"context"
	"errors"
)

type key int

const (
	cxtKey key = 0xACE0BA5E
)

func NewContext(ctx context.Context, logger *Client) context.Context {
	return context.WithValue(ctx, cxtKey, logger)
}

func FromContext(ctx context.Context) (*Client, error) {
	val := ctx.Value(cxtKey)
	if val == nil {
		return nil, errors.New("context does not contain a logger client")
	}
	logger, ok := val.(*Client)
	if !ok {
		return nil, errors.New("context value is not a logger client")
	}
	return logger, nil
}
