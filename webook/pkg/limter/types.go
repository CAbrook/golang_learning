package limter

import "context"

type Limiter interface {
	Limit(ctx context.Context, key string) (bool, error)
}
