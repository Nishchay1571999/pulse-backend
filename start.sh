export MINIO_ENDPOINT=localhost:9000
export MINIO_ROOT_USER=minioadmin
export MINIO_ROOT_PASSWORD=minioadmin
export REDIS_ADDR=localhost:6379
export QDRANT_URL=http://127.0.0.1:6333

# Run the server:
go run ./cmd/server
