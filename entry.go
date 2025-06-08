package logit

import "context"

type Entry struct {
	level Level
	tags  *tags
	msg   string
	attr  map[string]any
	err   error
}

func (e *Entry) Attribute(attr map[string]any) *Entry {
	e.attr = attr
	return e
}

func (e *Entry) Error(err error) *Entry {
	e.err = err
	return e
}
func (e *Entry) Log(ctx context.Context) error {
	return nil
}
