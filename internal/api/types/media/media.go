package media_models

import (
	"time"
)

type Post struct {
	PostId      int       `json:"postId" db:"post_id"`
	BlobName    string    `json:"blobName" db:"blob_name"`
	ContentType string    `json:"contentType" db:"content_type"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	Restricted  bool      `json:"restricted" db:"restricted"`
}
