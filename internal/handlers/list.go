package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/Nishchay1571999/pulse/backend/internal/storage"
)

type ListHandler struct {
    Storage storage.Storage
    DefaultBucket string
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    bucket := r.URL.Query().Get("bucket")
    if bucket == "" {
        bucket = h.DefaultBucket
    }
    prefix := r.URL.Query().Get("prefix")

    exists, err := h.Storage.BucketExists(r.Context(), bucket)
    if err != nil {
        log.Printf("bucket exists check: %v", err)
        http.Error(w, "failed to check bucket", http.StatusInternalServerError)
        return
    }
    if !exists {
        http.Error(w, "bucket not found", http.StatusNotFound)
        return
    }

    files, err := h.Storage.ListObjects(r.Context(), bucket, prefix)
    if err != nil {
        log.Printf("list objects: %v", err)
        http.Error(w, "list failed", http.StatusInternalServerError)
        return
    }

    if err := json.NewEncoder(w).Encode(files); err != nil {
        log.Printf("encode: %v", err)
        http.Error(w, "encode failed", http.StatusInternalServerError)
    }
}
