package main

import (
	"context"

	"github.com/g8rswimmer/go-logit"
)

func main() {
	logger := logit.NewClient()
	logger.AddTag("some_tag", "this is a tag").
		AddTag("number_tag", 1).
		Info("This is for some information").
		WithAttribute("attr1", "some attribute").
		WithAttribute("attr2", 1233).
		WithAttribute("attr3", map[string]any{
			"sub_attr3": "hi",
		}).Log(context.Background())
}
