package storage

import (
    "context"
    "fmt"
    "io"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStore struct {
    client *minio.Client
}

func NewMinioStore(endpoint, user, pass string, secure bool) (Storage, error) {
    if user == "" || pass == "" {
        return nil, fmt.Errorf("minio credentials missing")
    }
    cli, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(user, pass, ""),
        Secure: secure,
    })
    if err != nil {
        return nil, err
    }
    return &minioStore{client: cli}, nil
}

func (m *minioStore) BucketExists(ctx context.Context, bucket string) (bool, error) {
    return m.client.BucketExists(ctx, bucket)
}

func (m *minioStore) MakeBucket(ctx context.Context, bucket string) error {
    return m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

func (m *minioStore) PutObject(ctx context.Context, bucket, object string, data ReadSeekSizer, contentType string) error {
    // ensure we seek to start
    if _, err := data.Seek(0, io.SeekStart); err != nil {
        return err
    }
    _, err := m.client.PutObject(ctx, bucket, object, data, data.Size(), minio.PutObjectOptions{
        ContentType: contentType,
    })
    return err
}

func (m *minioStore) ListObjects(ctx context.Context, bucket, prefix string) ([]string, error) {
    opts := minio.ListObjectsOptions{
        Prefix:    prefix,
        Recursive: true,
    }
    var out []string
    for obj := range m.client.ListObjects(ctx, bucket, opts) {
        if obj.Err != nil {
            // log and skip; caller will get partial results
            continue
        }
        out = append(out, obj.Key)
    }
    return out, nil
}
