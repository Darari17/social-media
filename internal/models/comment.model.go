package models

import (
	"time"
)

type Comment struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	PostID    int        `db:"post_id"`
	Content   string     `db:"content"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
