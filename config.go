package logit

import (
	"io"
	"os"
	"time"
)

type Option func(cfg *config)

type config struct {
	levelConverter  LevelCoversion
	writer          io.Writer
	timeLayout      string
	timeStampLayout string
	timeStampField  string
	messageField    string
	levelField      string
}

var defaultConfig = config{
	levelConverter:  defaultLevelString,
	writer:          os.Stdout,
	timeLayout:      time.RFC3339,
	timeStampLayout: time.RFC3339,
	timeStampField:  "timestamp",
	messageField:    "message",
	levelField:      "level",
}

func SetDefaultConfiguration(opt Option, opts ...Option) {
	opts = append(opts, opt)
	for _, o := range opts {
		o(&defaultConfig)
	}
}

func WithLevelConverter(lc LevelCoversion) Option {
	return func(cfg *config) {
		cfg.levelConverter = lc
	}
}

func WithWriter(w io.Writer) Option {
	return func(cfg *config) {
		cfg.writer = w
	}
}

func WithTimeLayout(layout string) Option {
	return func(cfg *config) {
		cfg.timeLayout = layout
	}
}

func WithTimestampLayout(layout string) Option {
	return func(cfg *config) {
		cfg.timeStampLayout = layout
	}
}

func WithTimestampFieldName(name string) Option {
	return func(cfg *config) {
		cfg.timeStampField = name
	}
}

func WithMessageFieldName(name string) Option {
	return func(cfg *config) {
		cfg.messageField = name
	}
}

func WithLevelFieldName(name string) Option {
	return func(cfg *config) {
		cfg.levelField = name
	}
}
