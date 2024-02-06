package util

import (
	"context"
)

type (
	ISomeUtil interface {
		DoSomething(ctx context.Context) (string, error)
	}
)
