package dtos

import (
	"time"
)

type CommentRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type CommentResponse struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	PostID    int        `json:"post_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
