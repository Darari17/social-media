package models

import (
	"time"
)

type Post struct {
	ID        int        `db:"id"`
	UserID    int        `db:"user_id"`
	Content   *string    `db:"content_text"`
	Image     *string    `db:"content_image"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
