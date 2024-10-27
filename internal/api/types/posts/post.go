package posts

import (
	"time"
)

type Post struct {
	PostId     int       `json:"post_id" db:"post_id"`
	Title      string    `json:"title" db:"title"`
	Content    string    `json:"content" db:"content"`
	UserId     int       `json:"userId" db:"user_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	Restricted bool      `json:"restricted" db:"restricted"`
}

type FrontendPostRequest struct {
	PostRequestBody `json:"postData"`
}
type PostRequestBody struct {
	Title      string `json:"title" db:"title"`
	Content    string `json:"content" db:"content"`
	Restricted bool   `json:"restricted" db:"restricted"`
}

type CreatedPost struct {
	PostId int `json:"post_id" db:"post_id"`
}
