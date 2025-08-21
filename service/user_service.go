package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	"github.com/wahyujatirestu/sahabat-kurban/utils/security"
)

type UserService interface {
	GetAll(ctx context.Context) ([]dto.UserResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role string) (*dto.UserResponse, error)
	CreateWithRole(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAll(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := dto.ToUserResponseList(users)
	return res, nil
}

func (s *userService) GetById(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	users, err := s.userRepo.FindById(ctx, id)
	if err != nil || users == nil {
		return nil, errors.New("User not found")
	}
	
	if strings.TrimSpace(req.Username) != "" {
		users.Username = req.Username
	}
	if strings.TrimSpace(req.Name) != "" {
		users.Name = req.Name
	}
	if strings.TrimSpace(req.Email) != "" {
		users.Email = req.Email
	}

	user := &model.User{
		ID: id,
		Username: req.Username,
		Name: req.Name,
		Email: req.Email,
		Role: "user",
	}

	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userService) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil || user == nil {
		return errors.New("User not found")
	}

	valid, err := security.VerifyPasswordHash(user.Password, req.OldPassword)
	if err != nil || !valid {
		return errors.New("Old password is incorrect")
	}

	newHashed, err := security.GeneratePasswordHash(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = newHashed
	user.Updated_At = time.Now()

	return s.userRepo.UpdatePassword(ctx, userID, newHashed)
}

func (s *userService) UpdateRole(ctx context.Context, userID uuid.UUID, role string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("User not found")
	}

	if strings.TrimSpace(role) != "" {
		user.Role = role
	}

	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)
	return &res, nil
}

func (s *userService) CreateWithRole(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	hashed, err := security.GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("Name is required")
	}
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("Username is required")
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("Email is required")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("Password is required")
	}
	if strings.TrimSpace(req.Role) == "" {
		return nil, errors.New("Role is required")
	}

	user := &model.User{
		ID:       uuid.New(),
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     req.Role,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	res := dto.ToUserResponse(user)

	return &res, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
