package repo

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/interfaces"
)

type UserRepo struct {
	dao *pgxpool.Pool
}

func NewUserRepository(dao *pgxpool.Pool) interfaces.UserRepoInterface {
	return &UserRepo{
		dao: dao,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, username, email, password)
		VALUES($1, $2, $3, $4)
	`
	_, err := r.dao.Exec(ctx, query, user.ID, user.UserName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, ID string) (*entity.User, error) {
	query := `
		SELECT
			id, 
			username, 
			email
		FROM
			users
		WHERE
			id = $1
	`
	var (
		id, username, email sql.NullString
	)

	err := r.dao.QueryRow(ctx, query, ID).Scan(
		&id,
		&username,
		&email,
	)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       id.String,
		UserName: username.String,
		Email:    email.String,
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT 
			id, 
			username, 
			email,
			password
		FROM 
			users
		WHERE 
			email = $1
	`
	var (
		id, username, dbEmail, password sql.NullString
	)

	err := r.dao.QueryRow(ctx, query, email).Scan(
		&id,
		&username,
		&dbEmail,
		&password,
	)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       id.String,
		UserName: username.String,
		Email:    dbEmail.String,
		Password: password.String,
	}

	return user, nil
}
