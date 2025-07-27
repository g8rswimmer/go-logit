package main

import (
	"context"
	"errors"

	"github.com/g8rswimmer/go-logit"
)

func main() {
	logger := logit.NewClient(logit.WithFormatter(&logit.FormatJSON{}))
	ctx := logit.NewContext(context.Background(), logger)

	addTags(ctx)

	logWithError(ctx)

	anotherLogger := logger.Copy()
	anotherLogger.AddTag("another_tag", "this is another tag")
	anotherCtx := logit.NewContext(context.Background(), anotherLogger)
	logInfo(anotherCtx)

	logInfoStruct(ctx)
}

func addTags(ctx context.Context) {
	logger, err := logit.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	logger.AddTag("some_tag", "this is a tag").AddTag("number_tag", 1)
}

func logWithError(ctx context.Context) {
	logger, err := logit.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	logger.Error("Some error has happened here").
		WithError(errors.New("some error message")).
		WithAttribute("attr1", "some attribute").
		Log(ctx)
}

func logInfo(ctx context.Context) {
	logger, err := logit.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	logger.Info("This is for some information").
		WithAttribute("attr2", 1233).
		WithAttribute("attr3", map[string]any{
			"sub_attr3": "hi",
		}).Log(ctx)
}

func logInfoStruct(ctx context.Context) {
	logger, err := logit.FromContext(ctx)
	if err != nil {
		panic(err)
	}

	s := struct {
		SomeAttr string `logit:"some_attr"`
		SomeNum  int    `logit:"some_number"`
		NoSee    string `logit:"no_see_me,obfuscate"`
		OmitThis string `logit:",omit"`
	}{
		SomeAttr: "this is an attribute",
		SomeNum:  1234,
		NoSee:    "this will not be seen",
		OmitThis: "this will not be logged",
	}

	logger.Info("Going to have a structure").
		WithAttribute("log_struct", s).Log(ctx)
}
