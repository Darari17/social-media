package dtos

import (
	"mime/multipart"
	"time"
)

type PostRequest struct {
	Content string                `form:"content"`
	Image   *multipart.FileHeader `form:"image"`
}

type PostUpdateRequest struct {
	Content *string               `form:"content"`
	Image   *multipart.FileHeader `form:"image"`
}

type PostResponse struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Content   *string    `json:"content"`
	Image     *string    `json:"image"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
