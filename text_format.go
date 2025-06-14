package logit

import "context"

type TextFormat struct{}

func (t *TextFormat) Format(ctx context.Context, e *Entry) error {
	return nil
}
