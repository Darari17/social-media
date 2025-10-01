package repos

import (
	"context"
	"fmt"

	"github.com/Darari17/social-media/internal/dtos"
	"github.com/Darari17/social-media/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) GetAllUsers(c context.Context) ([]dtos.UserResponse, error) {
	query := "select id, name, email, avatar, bio, created_at, updated_at from users"

	rows, err := ur.db.Query(c, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dtos.UserResponse
	for rows.Next() {
		var user dtos.UserResponse
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Avatar, &user.Bio, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepo) GetUserByID(c context.Context, userId int) (*models.User, error) {
	query := "select id, name, email, avatar, bio, created_at, updated_at from users where id = $1"

	var user models.User

	if err := ur.db.QueryRow(c, query, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Avatar, &user.Bio, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) UpdateUser(c context.Context, user *models.User) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	i := 1

	if user.Name != nil {
		query += fmt.Sprintf("name = $%d,", i)
		args = append(args, *user.Name)
		i++
	}
	if user.Avatar != nil {
		query += fmt.Sprintf("avatar = $%d,", i)
		args = append(args, *user.Avatar)
		i++
	}
	if user.Bio != nil {
		query += fmt.Sprintf("bio = $%d,", i)
		args = append(args, *user.Bio)
		i++
	}

	if len(args) == 0 {
		return nil
	}
	query = query[:len(query)-1]
	query += fmt.Sprintf(", updated_at = now() WHERE id = $%d", i)
	args = append(args, user.ID)

	_, err := ur.db.Exec(c, query, args...)
	return err
}
