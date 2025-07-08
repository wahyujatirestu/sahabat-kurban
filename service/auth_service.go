package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/config"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	utilrepo "github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	"github.com/wahyujatirestu/sahabat-kurban/utils/security"
	"github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	cfg 		*config.Config
	userRepo 	repository.UserRepository
	rtRepo		utilrepo.RefreshTokenRepository
	jwtService 	service.JWTService
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, rtRepo utilrepo.RefreshTokenRepository, jwtService service.JWTService) AuthService {
	return &authService{
		cfg: cfg,
		userRepo: userRepo,
		rtRepo: rtRepo,
		jwtService: jwtService,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.AuthResponse, error) {
	existing, _ := s.userRepo.FindByEmailOrUsername(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("Email already used")
	}

	existing, _ = s.userRepo.FindByEmailOrUsername(ctx, req.Username)
	if existing != nil {
		return nil, errors.New("Username already used")
	}

	hashed, err := security.GeneratePasswordHash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID: uuid.New(),
		Username: req.Username,
		Name: req.Name,
		Email: req.Email,
		Password: hashed,
		Role: "user",
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil || user == nil {
		return nil, errors.New("Invalid credentials")
	}

	valid, err := security.VerifyPasswordHash(user.Password, req.Password)
	if err != nil || !valid {
		return nil, errors.New("Invalid credentials")
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	rt, err := s.rtRepo.FindByToken(ctx, refreshToken)

	if err != nil || rt == nil || rt.Revoked || rt.Expires_At.Before(time.Now()) {
		return nil, errors.New("Invalid or expired refresh token")
	}

	user, err := s.userRepo.FindById(ctx, rt.UserID)
	if err != nil || user == nil {
		return nil, errors.New("User not found")
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return  s.jwtService.RevokeRefreshTokenByToken(ctx, refreshToken)
}

func (s *authService) generateAuthResponse(ctx context.Context, user *model.User) (*dto.AuthResponse, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		ID:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}
