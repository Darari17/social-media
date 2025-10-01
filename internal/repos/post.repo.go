package repos

import (
	"context"
	"fmt"
	"strings"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) *PostRepo {
	return &PostRepo{db: db}
}

func (pr *PostRepo) CreatePost(c context.Context, post *models.Post) error {
	query := `INSERT INTO posts (user_id, content_text, content_image, created_at)
			  VALUES ($1, $2, $3, now()) returning id`
	err := pr.db.QueryRow(c, query, post.UserID, post.Content, post.Image).Scan(&post.ID)
	return err
}

func (pr *PostRepo) GetAllPosts(c context.Context) ([]dtos.PostResponse, error) {
	query := `SELECT id, user_id, content_text, content_image, created_at, updated_at, deleted_at 
	          FROM posts 
	          WHERE deleted_at IS NULL
	          ORDER BY COALESCE(updated_at, created_at) DESC
`
	rows, err := pr.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []dtos.PostResponse
	for rows.Next() {
		var p dtos.PostResponse
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Image, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (pr *PostRepo) GetPostsByUser(c context.Context, userId int) ([]dtos.PostResponse, error) {
	query := `SELECT id, user_id, content_text, content_image, created_at, updated_at, deleted_at 
	          FROM posts 
	          WHERE user_id=$1 AND deleted_at IS NULL
	          ORDER BY created_at DESC`
	rows, err := pr.db.Query(c, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []dtos.PostResponse
	for rows.Next() {
		var p dtos.PostResponse
		if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Image, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (pr *PostRepo) GetPostByID(c context.Context, id int) (*dtos.PostResponse, error) {
	query := `SELECT id, user_id, content_text, content_image, created_at, updated_at, deleted_at 
	          FROM posts 
	          WHERE id=$1 AND deleted_at IS NULL`
	var p dtos.PostResponse
	if err := pr.db.QueryRow(c, query, id).Scan(&p.ID, &p.UserID, &p.Content, &p.Image, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

func (pr *PostRepo) UpdatePost(c context.Context, post *models.Post) error {
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if post.Content != nil {
		setClauses = append(setClauses, fmt.Sprintf("content_text=$%d", argID))
		args = append(args, post.Content)
		argID++
	}

	if post.Image != nil {
		setClauses = append(setClauses, fmt.Sprintf("content_image=$%d", argID))
		args = append(args, post.Image)
		argID++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at=now()")

	query := fmt.Sprintf(`UPDATE posts SET %s WHERE id=$%d`,
		strings.Join(setClauses, ", "), argID)
	args = append(args, post.ID)

	_, err := pr.db.Exec(c, query, args...)
	return err
}

func (pr *PostRepo) DeletePost(c context.Context, postId int) error {
	query := `UPDATE posts SET deleted_at = now() WHERE id=$1`
	_, err := pr.db.Exec(c, query, postId)
	return err
}
