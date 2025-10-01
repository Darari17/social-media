package repos

import (
	"context"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

func (cr *CommentRepo) CreateComment(c context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (user_id, post_id, content, created_at)
	          VALUES ($1, $2, $3, now())`
	_, err := cr.db.Exec(c, query, comment.UserID, comment.PostID, comment.Content)
	return err
}

func (cr *CommentRepo) GetCommentsByPost(c context.Context, postId int) ([]dtos.CommentResponse, error) {
	query := `SELECT id, user_id, post_id, content, created_at, updated_at 
	          FROM comments 
	          WHERE post_id=$1
	          ORDER BY created_at ASC`

	rows, err := cr.db.Query(c, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []dtos.CommentResponse
	for rows.Next() {
		var cm dtos.CommentResponse
		if err := rows.Scan(&cm.ID, &cm.UserID, &cm.PostID, &cm.Content, &cm.CreatedAt, &cm.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, cm)
	}
	return comments, nil
}

func (cr *CommentRepo) UpdateComment(c context.Context, commentId int, content string) error {
	query := `UPDATE comments SET content=$1, updated_at=now() WHERE id=$2`
	_, err := cr.db.Exec(c, query, content, commentId)
	return err
}

func (cr *CommentRepo) DeleteComment(c context.Context, commentId int) error {
	query := `DELETE FROM comments WHERE id=$1`
	_, err := cr.db.Exec(c, query, commentId)
	return err
}
