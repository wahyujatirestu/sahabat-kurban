package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/config"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/repository"
	emailModel "github.com/wahyujatirestu/sahabat-kurban/utils/model"
	utilrepo "github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	"github.com/wahyujatirestu/sahabat-kurban/utils/security"
	"github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenOnlyResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerification(ctx context.Context, email string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	cfg 		*config.Config
	userRepo 	repository.UserRepository
	rtRepo		utilrepo.RefreshTokenRepository
	emailRepo	utilrepo.EmailVerificationRepository
	resetRepo 	utilrepo.ResetPasswordRepository
	jwtService 	service.JWTService
	emailService service.EmailService
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, rtRepo utilrepo.RefreshTokenRepository, emailRepo utilrepo.EmailVerificationRepository, resetRepo utilrepo.ResetPasswordRepository, jwtService service.JWTService, emailService service.EmailService) AuthService {
	return &authService{
		cfg: cfg,
		userRepo: userRepo,
		rtRepo: rtRepo,
		emailRepo: emailRepo,
		resetRepo: resetRepo,
		jwtService: jwtService,
		emailService: emailService,
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
		IsVerified: false,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	verifToken := uuid.New().String()
	tokenModel := &emailModel.EmailVerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verifToken,
		Expires_At: time.Now().Add(5 * time.Minute),
		Created_At: time.Now(),
	}
	
	if err := s.emailRepo.Save(ctx, tokenModel); err != nil {
		return nil, err
	}

	log.Println("Link verifikasi:", fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.cfg.AppBaseURL, verifToken))


	if err := s.emailService.SendVerificationEmail(user.Email, user.Name, verifToken); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByEmailOrUsername(ctx, req.Identifier)
	if err != nil || user == nil {
		return nil, errors.New("Invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("Email is not verified")
	}

	valid, err := security.VerifyPasswordHash(user.Password, req.Password)
	if err != nil || !valid {
		return nil, errors.New("Invalid credentials")
	}

	return s.generateAuthResponse(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenOnlyResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("Invalid or expired refresh token")
	}

	rt, err := s.rtRepo.FindByToken(ctx, refreshToken)
	if err != nil || rt == nil || rt.Revoked || rt.Expires_At.Before(time.Now()) {
		return nil, errors.New("Refresh token invalid or revoked")
	}

	userID, err := uuid.Parse(claims.UserId)
	if err != nil {
		return nil, errors.New("Invalid user ID")
	}

	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("User not found")
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(ctx, userID)

	return &dto.TokenOnlyResponse{
		AccessToken: accessToken,
		RefreshToken: newRefreshToken, 
	}, nil
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
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	t, err := s.emailRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("Invalid token")
	}

	if t.Expires_At.Before(time.Now()) {
		return errors.New("Token expired")
	}

	user, err := s.userRepo.FindById(ctx, t.UserID)
	if err != nil || user == nil {
		return errors.New("User not found")
	}

	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.emailRepo.DeleteByUserID(ctx, user.ID)

	return nil
}

func (s *authService) ResendVerification(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("User not found")
	}

	if user.IsVerified {
		return errors.New("Email is already verified")
	}

	_ = s.emailRepo.DeleteByUserID(ctx, user.ID)

	token := uuid.New().String()
	verif := &emailModel.EmailVerificationToken{
		ID: uuid.New(),
		UserID: user.ID,
		Token: token,
		Expires_At: time.Now().Add(24 * time.Hour),
		Created_At: time.Now(),
	}
	if err := s.emailRepo.Save(ctx, verif); err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(user.Email, user.Name, token)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("User tidak ditemukan")
	}

	if !user.IsVerified {
		return errors.New("Email belum diverifikasi")
	}

	_ = s.resetRepo.DeleteByUserID(ctx, user.ID)

	token := uuid.New().String()
	reset := &emailModel.ResetPasswordToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := s.resetRepo.Save(ctx, reset); err != nil {
		return err
	}

	return s.emailService.SendResetPasswordEmail(user.Email, user.Name, token)
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := s.resetRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("Token tidak valid")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("Token telah kadaluarsa")
	}

	hashedPassword, err := security.GeneratePasswordHash(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, resetToken.UserID, hashedPassword); err != nil {
		return err
	}

	return s.resetRepo.DeleteByUserID(ctx, resetToken.UserID)
}
