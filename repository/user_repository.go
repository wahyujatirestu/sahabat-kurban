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
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
	FindById(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, newHashed string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, username, name, email, password, role, is_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, user.ID, user.Username, user.Name, user.Email, user.Password, user.Role, user.IsVerified, user.Created_At, user.Updated_At)
	return err
}

func (r *userRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, username, name, email, password, role, is_verified, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user []*model.User
	for rows.Next(){
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsVerified, &u.Created_At, &u.Updated_At); err != nil {
			return nil, err
		}

		user = append(user, &u)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, is_verified, created_at, updated_at FROM users WHERE email=$1`, email)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsVerified, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err	
}

func (r *userRepository) FindByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, is_verified, created_at, updated_at FROM users WHERE email=$1 OR username=$1`, identifier)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsVerified, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err
}

func (r *userRepository) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, username, name, email, password, role, is_verified, created_at, updated_at FROM users WHERE id=$1`, id)

	var u model.User
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsVerified, &u.Created_At, &u.Updated_At)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &u, err	
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET username=$1, name=$2, email=$3, role=$4, is_verified=$5 WHERE id=$6`, user.Username, user.Name, user.Email, user.Role, user.IsVerified, user.ID)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newHashed string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET password=$1 WHERE id=$2`, newHashed, userID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("User not found")
	}

	return nil
}