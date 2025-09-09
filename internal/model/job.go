package model

type Job struct {
    ID     string `json:"id"`
    Bucket string `json:"bucket"`
    Object string `json:"object"`
}
