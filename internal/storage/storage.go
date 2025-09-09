package storage

import "context"
import "io"

type ReadSeekSizer interface {
    io.ReadSeeker
    Size() int64
}

// Storage abstracts object storage operations
type Storage interface {
    BucketExists(ctx context.Context, bucket string) (bool, error)
    MakeBucket(ctx context.Context, bucket string) error
    PutObject(ctx context.Context, bucket, object string, data ReadSeekSizer, contentType string) error
    ListObjects(ctx context.Context, bucket, prefix string) ([]string, error)
}
