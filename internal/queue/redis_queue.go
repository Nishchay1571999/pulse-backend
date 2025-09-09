package queue

import (
    "context"
    "github.com/go-redis/redis/v8"
)

type redisQueue struct {
    cli *redis.Client
    key string
}

func NewRedisQueue(addr string) (Queue, error) {
    cli := redis.NewClient(&redis.Options{Addr: addr})
    if _, err := cli.Ping(context.Background()).Result(); err != nil {
        return nil, err
    }
    return &redisQueue{cli: cli, key: "jobs"}, nil
}

func (r *redisQueue) Push(ctx context.Context, key string, payload []byte) error {
    k := r.key
    if key != "" {
        k = key
    }
    return r.cli.LPush(ctx, k, payload).Err()
}
