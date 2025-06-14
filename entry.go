package logit

import "context"

type Entry struct {
	level Level
	cfg   *Config
	tags  *Tags
	msg   string
	attr  map[string]any
	err   error
}

func (e Entry) Level() Level {
	return e.level
}

func (e Entry) Config() *Config {
	return e.cfg
}

func (e Entry) Tags() *Tags {
	return e.tags
}

func (e Entry) Message() string {
	return e.msg
}

func (e Entry) Atrributes() map[string]any {
	return e.attr
}

func (e Entry) Error() error {
	return e.err
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
