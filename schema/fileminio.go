package schema

import "time"

type FileMINIO struct {
	Bucket       string    `json:"bucket"`
	Filename     string    `json:"filename"`
	LastModified time.Time `json:"last_modified"`
}
