package dtos

import (
	"mime/multipart"
	"time"
)

type UserRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Name   *string               `json:"name" form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
	Bio    *string               `json:"bio" form:"bio"`
}

type UserTokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID        int        `json:"id"`
	Name      *string    `json:"name"`
	Email     *string    `json:"email"`
	Avatar    *string    `json:"avatar"`
	Bio       *string    `json:"bio"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
