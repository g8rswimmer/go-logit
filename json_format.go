package logit

import (
	"context"
)

type JSONFormat struct{}

func (j *JSONFormat) Format(ctx context.Context, e *Entry) error {
	return nil
}
