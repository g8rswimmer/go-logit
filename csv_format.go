package logit

import "context"

type CSVFormat struct{}

func (c *CSVFormat) Format(ctx context.Context, e *Entry) error {
	return nil
}
