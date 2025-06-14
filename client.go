package logit

import "fmt"

type Client struct {
	tags *Tags
	cfg  *Config
}

func NewClient(opts ...Option) *Client {
	cfg := Config(defaultConfig)
	for _, o := range opts {
		o(&cfg)
	}

	return &Client{
		cfg: &cfg,
		tags: &Tags{
			entries: map[string]any{},
		},
	}
}

func (c *Client) Copy() *Client {
	return &Client{
		tags: c.tags.copy(),
	}
}

func (c *Client) AddTag(name string, value any) *Client {
	c.tags.add(name, value)
	return c
}

func (c *Client) Trace(msg string, args ...any) *Entry {
	return c.Entry(LevelTrace, msg, args...)
}

func (c *Client) Debug(msg string, args ...any) *Entry {
	return c.Entry(LevelDebug, msg, args...)
}

func (c *Client) Info(msg string, args ...any) *Entry {
	return c.Entry(LevelInfo, msg, args...)
}

func (c *Client) Warn(msg string, args ...any) *Entry {
	return c.Entry(LevelWarn, msg, args...)
}

func (c *Client) Error(msg string, args ...any) *Entry {
	return c.Entry(LevelError, msg, args...)
}

func (c *Client) Critical(msg string, args ...any) *Entry {
	return c.Entry(LevelCritical, msg, args...)
}

func (c *Client) Emergency(msg string, args ...any) *Entry {
	return c.Entry(LevelEmergency, msg, args...)
}

func (c *Client) Fatal(msg string, args ...any) *Entry {
	return c.Entry(LevelFatal, msg, args...)
}

func (c *Client) Entry(level Level, msg string, args ...any) *Entry {
	return &Entry{
		level: level,
		cfg:   c.cfg,
		tags:  c.tags.copy(),
		msg:   fmt.Sprintf(msg, args...),
		attr:  map[string]any{},
	}
}
