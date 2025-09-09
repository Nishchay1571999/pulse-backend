package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/Nishchay1571999/pulse/backend/internal/model"
    "github.com/Nishchay1571999/pulse/backend/internal/queue"
    "github.com/Nishchay1571999/pulse/backend/internal/storage"
    "github.com/Nishchay1571999/pulse/backend/internal/util"
)

type UploadHandler struct {
    Storage storage.Storage
    Queue   queue.Queue
    DefaultBucket string
}

// helper to copy a multipart.File to temp file and return path
func saveMultipartToTemp(file multipart.File) (string, error) {
    tmpPath := filepath.Join(os.TempDir(), "pulse_"+uuid.New().String())
    out, err := os.Create(tmpPath)
    if err != nil {
        return "", err
    }
    defer out.Close()
    if _, err := io.Copy(out, file); err != nil {
        os.Remove(tmpPath)
        return "", err
    }
    return tmpPath, nil
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // limit handled by server if desired
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        http.Error(w, "parse form: "+err.Error(), http.StatusBadRequest)
        return
    }

    f, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "missing file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer f.Close()

    bucket := h.DefaultBucket
    exists, err := h.Storage.BucketExists(r.Context(), bucket)
    if err != nil {
        log.Printf("bucket exists error: %v", err)
        http.Error(w, "storage error", http.StatusInternalServerError)
        return
    }
    if !exists {
        if err := h.Storage.MakeBucket(r.Context(), bucket); err != nil {
            log.Printf("make bucket error: %v", err)
            http.Error(w, "storage error", http.StatusInternalServerError)
            return
        }
    }

    tmpPath, err := saveMultipartToTemp(f)
    if err != nil {
        http.Error(w, "tmp write: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer os.Remove(tmpPath)

    fw, err := util.NewFileWithSize(tmpPath)
    if err != nil {
        http.Error(w, "tmp open: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer fw.Close()

    objectName := fmt.Sprintf("%s-%s", uuid.New().String(), header.Filename)
    if err := h.Storage.PutObject(r.Context(), bucket, objectName, fw, header.Header.Get("Content-Type")); err != nil {
        log.Printf("put object error: %v", err)
        http.Error(w, "storage put error", http.StatusInternalServerError)
        return
    }

    job := model.Job{ID: uuid.New().String(), Bucket: bucket, Object: objectName}
    b, _ := json.Marshal(job)
    if err := h.Queue.Push(r.Context(), "", b); err != nil {
        log.Printf("queue push error: %v", err)
        http.Error(w, "queue error", http.StatusInternalServerError)
        return
    }

    resp := map[string]string{"job_id": job.ID, "object": objectName}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
