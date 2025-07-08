package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
)


type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, username, name, email, password, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, user.ID, user.Username, user.Name, user.Email, user.Password, user.Role, user.Created_At, user.Updated_At)
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, created_at, updated_at FROM users WHERE email=$1`, email)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err	
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, created_at, updated_at FROM users WHERE email=$1 OR username=$1`, identifier)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err
}

func (r *userRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, created_at, updated_at FROM users WHERE id=$1`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err	
}