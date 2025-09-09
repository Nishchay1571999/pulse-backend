package queue

import (
    "context"
)

type Queue interface {
    Push(ctx context.Context, key string, payload []byte) error
}
