package logit

import "context"

type Formatter interface {
	Format(ctx context.Context, e *Entry) error
}
