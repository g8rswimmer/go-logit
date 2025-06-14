package logit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type TextFormat struct{}

func (t *TextFormat) formatTags(e *Entry) ([]string, error) {
	encodedTags, err := encode(e.tags.Retrieve())
	if err != nil {
		return []string{}, err
	}
	tags, ok := encodedTags.(map[string]any)
	if !ok {
		return []string{}, errors.New("tags are not a map")
	}
	tagEntries := make([]string, len(tags))
	idx := 0
	for k, v := range tags {
		enc, err := json.Marshal(v)
		if err != nil {
			enc = []byte("error: unable to marshal tag[" + k + "]")
		}
		tagEntries[idx] = k + ":" + string(enc)
		idx++
	}
	return tagEntries, nil
}

func (t *TextFormat) formatAttributes(e *Entry) ([]string, error) {
	encodedAttributes, err := encode(e.attr)
	if err != nil {
		return []string{}, err
	}
	attributes, ok := encodedAttributes.(map[string]any)
	if !ok {
		return []string{}, errors.New("tags are not a map")
	}
	attrEntries := make([]string, len(attributes))
	idx := 0
	for k, v := range attributes {
		enc, err := json.Marshal(v)
		if err != nil {
			enc = []byte("error: unable to marshal tag[" + k + "]")
		}
		attrEntries[idx] = k + ":" + string(enc)
		idx++
	}
	return attrEntries, nil
}
func (t *TextFormat) Format(ctx context.Context, e *Entry) error {
	textEntry := fmt.Sprintf("%s %s %s:%s",
		time.Now().Format(e.cfg.timeStampLayout),
		e.cfg.levelConverter[e.level],
		e.cfg.messageField,
		e.msg)
	textTags, err := t.formatTags(e)
	switch {
	case err != nil:
		textEntry += fmt.Sprintf(" %s:[%s]", e.cfg.tagsField, err.Error())
	case len(textTags) > 0:
		textEntry += fmt.Sprintf(" %s:[%s]", e.cfg.tagsField, strings.Join(textTags, ", "))
	default:
	}
	textAttrs, err := t.formatAttributes(e)
	switch {
	case err != nil:
		textEntry += fmt.Sprintf(" %s:[%s]", e.cfg.attrField, err.Error())
	case len(textTags) > 0:
		textEntry += fmt.Sprintf(" %s:[%s]", e.cfg.attrField, strings.Join(textAttrs, ", "))
	default:
	}
	textEntry += "\n"
	_, err = e.cfg.writer.Write([]byte(textEntry))
	return err
}
