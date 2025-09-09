# ensure module mode
go env -w GO111MODULE=on

# fetch packages without specifying versions (Go will pick)
go get github.com/gorilla/mux
go get github.com/minio/minio-go/v7
go get github.com/go-redis/redis/v8
go get github.com/google/uuid

# now tidy to populate go.sum and prune unused ones
go mod tidy -v
