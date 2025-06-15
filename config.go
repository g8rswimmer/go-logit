package logit

import (
	"io"
	"os"
	"time"
)

type Config struct {
	levelConverter  LevelConversion
	formatter       Formatter
	writer          io.Writer
	timeLayout      string
	timeStampLayout string
	timeStampField  string
	messageField    string
	levelField      string
	tagsField       string
	attrField       string
	errField        string
}

var defaultConfig = Config{
	levelConverter:  defaultLevelString,
	writer:          os.Stdout,
	formatter:       &FormatText{},
	timeLayout:      time.RFC3339,
	timeStampLayout: time.RFC3339,
	timeStampField:  "timestamp",
	messageField:    "message",
	levelField:      "level",
	tagsField:       "tags",
	attrField:       "attributes",
	errField:        "error",
}

func SetDefaultConfiguration(opt Option, opts ...Option) {
	opts = append(opts, opt)
	for _, o := range opts {
		o(&defaultConfig)
	}
}

func (c Config) LevelConverter() LevelConversion {
	return c.levelConverter
}

func (c Config) Writer() io.Writer {
	return c.writer
}

func (c Config) TimeLayout() string {
	return c.timeLayout
}

func (c Config) TimeStampLayout() string {
	return c.timeStampLayout
}

func (c Config) TimeStampField() string {
	return c.timeStampField
}

func (c Config) MessageField() string {
	return c.messageField
}

func (c Config) LevelField() string {
	return c.levelField
}

func (c Config) TagsField() string {
	return c.tagsField
}

func (c Config) AttributesField() string {
	return c.attrField
}

func (c Config) ErrorField() string {
	return c.errField
}
