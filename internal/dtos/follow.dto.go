package dtos

import "time"

type FollowRequest struct {
	FollowingID int `json:"following_id" binding:"required"`
}

type FollowResponse struct {
	ID          int       `json:"id"`
	FollowerID  int       `json:"follower_id"`
	FollowingID int       `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}
