package models

import (
	"time"
)

type Like struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	PostID    int       `db:"post_id"`
	CreatedAt time.Time `db:"created_at"`
}
