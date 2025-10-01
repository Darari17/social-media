package repos

import (
	"context"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FollowRepo struct {
	db *pgxpool.Pool
}

func NewFollowRepo(db *pgxpool.Pool) *FollowRepo {
	return &FollowRepo{db: db}
}

func (fr *FollowRepo) FollowUser(c context.Context, followerId, followingId int) (*dtos.FollowResponse, error) {
	query := `
		INSERT INTO follows (follower_id, following_id, created_at)
		VALUES ($1, $2, now())
		ON CONFLICT (follower_id, following_id) DO NOTHING
		RETURNING id, follower_id, following_id, created_at
	`

	var res dtos.FollowResponse
	err := fr.db.QueryRow(c, query, followerId, followingId).
		Scan(&res.ID, &res.FollowerID, &res.FollowingID, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (fr *FollowRepo) UnfollowUser(c context.Context, followerId, followingId int) (int64, error) {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`
	cmdTag, err := fr.db.Exec(c, query, followerId, followingId)
	if err != nil {
		return 0, err
	}
	return cmdTag.RowsAffected(), nil
}

func (fr *FollowRepo) GetFollowers(c context.Context, userId int) ([]dtos.UserResponse, error) {
	query := `
		SELECT u.id, u.name, u.email, u.avatar, u.bio, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1
	`
	rows, err := fr.db.Query(c, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []dtos.UserResponse
	for rows.Next() {
		var u dtos.UserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Avatar, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, u)
	}
	return followers, nil
}

func (fr *FollowRepo) GetFollowing(c context.Context, userId int) ([]dtos.UserResponse, error) {
	query := `
		SELECT u.id, u.name, u.email, u.avatar, u.bio, u.created_at, u.updated_at
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1
	`
	rows, err := fr.db.Query(c, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []dtos.UserResponse
	for rows.Next() {
		var u dtos.UserResponse
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Avatar, &u.Bio, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		following = append(following, u)
	}
	return following, nil
}
