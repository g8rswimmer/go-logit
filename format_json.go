package logit

import (
	"context"
	"encoding/json"
	"time"
)

type FormatJSON struct{}

func (j *FormatJSON) Format(ctx context.Context, e *Entry) error {
	entry := map[string]any{
		e.cfg.messageField:   e.msg,
		e.cfg.timeStampField: time.Now().Format(e.cfg.timeStampLayout),
		e.cfg.levelField:     e.cfg.levelConverter[e.level],
	}
	if e.err != nil {
		entry[e.cfg.errField] = e.err.Error()
	}
	if t := e.tags.Retrieve(); len(t) > 0 {
		enc, err := encode(t)
		switch {
		case err != nil:
			enc = err.Error()
		default:
		}
		entry[e.cfg.tagsField] = enc
	}
	if len(e.attr) > 0 {
		enc, err := encode(e.attr)
		switch {
		case err != nil:
			enc = err.Error()
		default:
		}
		entry[e.cfg.attrField] = enc
	}
	return json.NewEncoder(e.cfg.writer).Encode(entry)
}
