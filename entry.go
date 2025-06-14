package logit

import "context"

type Entry struct {
	level Level
	cfg   *config
	tags  *tags
	msg   string
	attr  map[string]any
	err   error
}

func (e *Entry) WithAttribute(attr map[string]any) *Entry {
	e.attr = attr
	return e
}

func (e *Entry) WithError(err error) *Entry {
	e.err = err
	return e
}
func (e *Entry) Log(ctx context.Context) error {
	return nil
}
