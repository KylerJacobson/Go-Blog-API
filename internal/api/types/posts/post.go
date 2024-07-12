package post_models

import (
	"time"
)

type Post struct {
	PostId      int       `json:"postId" db:"post_id"`
	Title       string    `json:"title" db:"title"` 
	Content     string    `json:"content" db:"content"`
	UserId      int       `json:"userId" db:"user_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Restricted  bool      `json:"restricted" db:"restricted"`
}
