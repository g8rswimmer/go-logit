package logit

import "io"

type Option func(cfg *Config)

func WithLevelConverter(lc LevelConversion) Option {
	return func(cfg *Config) {
		cfg.levelConverter = lc
	}
}

func WithWriter(w io.Writer) Option {
	return func(cfg *Config) {
		cfg.writer = w
	}
}

func WithTimeLayout(layout string) Option {
	return func(cfg *Config) {
		cfg.timeLayout = layout
	}
}

func WithTimestampLayout(layout string) Option {
	return func(cfg *Config) {
		cfg.timeStampLayout = layout
	}
}

func WithTimestampFieldName(name string) Option {
	return func(cfg *Config) {
		cfg.timeStampField = name
	}
}

func WithMessageFieldName(name string) Option {
	return func(cfg *Config) {
		cfg.messageField = name
	}
}

func WithLevelFieldName(name string) Option {
	return func(cfg *Config) {
		cfg.levelField = name
	}
}

func WithTagsFieldName(name string) Option {
	return func(cfg *Config) {
		cfg.tagsField = name
	}
}

func WithAttributesFieldName(name string) Option {
	return func(cfg *Config) {
		cfg.attrField = name
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(cfg *Config) {
		cfg.formatter = formatter
	}
}
