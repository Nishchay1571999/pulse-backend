package util

import (
    "io"
    "os"
)

type FileWithSize struct {
    *os.File
    sz int64
}

func NewFileWithSize(path string) (*FileWithSize, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    fi, err := f.Stat()
    if err != nil {
        f.Close()
        return nil, err
    }
    return &FileWithSize{File: f, sz: fi.Size()}, nil
}

func (f *FileWithSize) Size() int64 {
    return f.sz
}

var _ io.ReadSeeker = (*FileWithSize)(nil)
