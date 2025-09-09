package config

import "os"

type Config struct {
    Addr         string
    MinioEndpoint string
    MinioUser    string
    MinioPass    string
    MinioSecure  bool
    RedisAddr    string
    DefaultBucket string
}

func LoadFromEnv() *Config {
    secure := false
    if os.Getenv("MINIO_SECURE") == "1" || os.Getenv("MINIO_SECURE") == "true" {
        secure = true
    }

    addr := os.Getenv("HTTP_ADDR")
    if addr == "" {
        addr = ":8080"
    }

    bucket := os.Getenv("DEFAULT_BUCKET")
    if bucket == "" {
        bucket = "videos"
    }

    return &Config{
        Addr:          addr,
        MinioEndpoint: envOr("MINIO_ENDPOINT", "localhost:9000"),
        MinioUser:     envOr("MINIO_ROOT_USER", ""),
        MinioPass:     envOr("MINIO_ROOT_PASSWORD", ""),
        MinioSecure:   secure,
        RedisAddr:     envOr("REDIS_ADDR", "localhost:6379"),
        DefaultBucket: bucket,
    }
}

func envOr(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}
