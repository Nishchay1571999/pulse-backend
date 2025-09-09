package main

import (
    "context"
    "log"

    "github.com/Nishchay1571999/pulse/backend/internal/config"
    "github.com/Nishchay1571999/pulse/backend/internal/queue"
    "github.com/Nishchay1571999/pulse/backend/internal/server"
    "github.com/Nishchay1571999/pulse/backend/internal/storage"
)

func main() {
    // load config from env
    cfg := config.LoadFromEnv()

    // init dependencies
    minioStore, err := storage.NewMinioStore(cfg.MinioEndpoint, cfg.MinioUser, cfg.MinioPass, cfg.MinioSecure)
    if err != nil {
        log.Fatalf("init minio: %v", err)
    }

    redisQ, err := queue.NewRedisQueue(cfg.RedisAddr)
    if err != nil {
        log.Fatalf("init redis: %v", err)
    }

    srv := server.New(&server.Dependencies{
        Storage: minioStore,
        Queue:   redisQ,
        Config:  cfg,
    })

    ctx := context.Background()
    if err := srv.Run(ctx); err != nil {
        log.Fatalf("server exited: %v", err)
    }
}
