package server

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/Nishchay1571999/pulse/backend/internal/config"
    "github.com/Nishchay1571999/pulse/backend/internal/handlers"
    "github.com/Nishchay1571999/pulse/backend/internal/queue"
    "github.com/Nishchay1571999/pulse/backend/internal/storage"
)

type Dependencies struct {
    Storage storage.Storage
    Queue   queue.Queue
    Config  *config.Config
}

type Server struct {
    deps *Dependencies
    srv  *http.Server
}

func New(deps *Dependencies) *Server {
    return &Server{deps: deps}
}

func (s *Server) Run(ctx context.Context) error {
    r := mux.NewRouter()
    // health
    r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"ok":true}`))
    }).Methods("GET")

    // handlers (inject dependencies)
    uploadH := &handlers.UploadHandler{
        Storage: s.deps.Storage,
        Queue:   s.deps.Queue,
        DefaultBucket: s.deps.Config.DefaultBucket,
    }
    listH := &handlers.ListHandler{
        Storage: s.deps.Storage,
        DefaultBucket: s.deps.Config.DefaultBucket,
    }

    r.Handle("/upload", uploadH).Methods("POST")
    r.Handle("/list", listH).Methods("GET")

    s.srv = &http.Server{
        Handler:      r,
        Addr:         s.deps.Config.Addr,
        WriteTimeout: 30 * time.Second,
        ReadTimeout:  30 * time.Second,
    }

    log.Printf("Pulse backend starting on %s", s.deps.Config.Addr)
    return s.srv.ListenAndServe()
}
