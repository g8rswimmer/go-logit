package logit

import (
	"reflect"
	"strings"
)

const (
	tagName      = "logit"
	omitTag      = "omit"
	obfuscateTag = "obfuscate"
)

type fieldTag struct {
	Name      string
	Omit      bool
	Obfuscate bool
}

func encodeFieldTag(field reflect.StructField) *fieldTag {
	ft := &fieldTag{
		Name: strings.ToLower(field.Name),
	}
	tagValue := field.Tag.Get(tagName)
	tagOptions := strings.Split(tagValue, ",")
	if len(tagOptions) == 0 {
		return ft
	}
	if len(tagOptions[0]) > 0 {
		ft.Name = tagOptions[0]
	}
	opts := map[string]any{}
	for i := 1; i < len(tagOptions); i++ {
		opts[tagOptions[i]] = nil
	}
	if _, has := opts[omitTag]; has {
		ft.Omit = true
	}
	if _, has := opts[obfuscateTag]; has {
		ft.Obfuscate = true
	}
	return ft
}
