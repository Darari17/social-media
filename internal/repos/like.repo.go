package repos

import (
	"context"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeRepo struct {
	db *pgxpool.Pool
}

func NewLikeRepo(db *pgxpool.Pool) *LikeRepo {
	return &LikeRepo{db: db}
}

func (lr *LikeRepo) CreateLike(c context.Context, like *models.Like) error {
	query := `INSERT INTO likes (user_id, post_id, created_at)
	          VALUES ($1, $2, now())`
	_, err := lr.db.Exec(c, query, like.UserID, like.PostID)
	return err
}

func (lr *LikeRepo) DeleteLike(c context.Context, userId, postId int) error {
	query := `DELETE FROM likes WHERE user_id=$1 AND post_id=$2`
	_, err := lr.db.Exec(c, query, userId, postId)
	return err
}

func (lr *LikeRepo) GetLikesByPost(c context.Context, postId int) ([]dtos.UserResponse, error) {
	query := `
		SELECT u.id, u.name, u.email, u.avatar, u.bio, u.created_at, u.updated_at
		FROM likes l
		JOIN users u ON l.user_id = u.id
		WHERE l.post_id = $1
	`

	rows, err := lr.db.Query(c, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dtos.UserResponse
	for rows.Next() {
		var user dtos.UserResponse
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Avatar,
			&user.Bio,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
