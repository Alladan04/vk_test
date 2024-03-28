package repo

import (
	"context"

	"github.com/Alladan04/vk_test/internal/models"
	"github.com/jackc/pgtype/pgxtype"
)

const (
	getUserByUsername = "SELECT id, username, password_hash, create_time, image_path FROM users WHERE username = $1; "
	addUser           = "INSERT INTO users(id, username, password_hash, create_time, image_path) VALUES ($1, $2, $3, $4, $5);"
)

type AuthRepo struct {
	db pgxtype.Querier
}

func NewAuthRepo(db pgxtype.Querier) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (repo *AuthRepo) AddUser(ctx context.Context, user models.User) error {
	_, err := repo.db.Exec(ctx, addUser, user.Id, user.Username, user.Password, user.CreateTime, user.Img)
	if err != nil {
		return err
	}
	return nil
}
func (repo *AuthRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {

	resultUser := models.User{Username: username}

	err := repo.db.QueryRow(ctx, getUserByUsername, username).Scan(
		&resultUser.Id,
		&resultUser.Username,
		&resultUser.Password,
		&resultUser.CreateTime,
		&resultUser.Img,
	)

	if err != nil {

		return models.User{}, err
	}

	return resultUser, nil
}
