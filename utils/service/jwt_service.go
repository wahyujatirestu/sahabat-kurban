package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/config"
	"github.com/wahyujatirestu/sahabat-kurban/utils/model"
	"github.com/wahyujatirestu/sahabat-kurban/utils/repository"
	userutils "github.com/wahyujatirestu/sahabat-kurban/repository"
)

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
	ValidateAccessToken(tokenStr string) (*model.JWTPayloadClaim, error)
	ValidateRefreshToken(ctx context.Context, tokenStr string) (*model.JWTPayloadClaim, error)
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	RevokeRefreshTokenByToken(ctx context.Context, token string) error
	DeleteExpiredToken(ctx context.Context) error
}

type jwtService struct {
	cfg		*config.Config
	refreshToken repository.RefreshTokenRepository
	userRepo userutils.UserRepository
}

func NewJWTServie(cfg *config.Config, repo repository.RefreshTokenRepository, userRepo userutils.UserRepository) JWTService {
	return &jwtService {
		cfg: cfg, refreshToken: repo, userRepo: userRepo,
	}
}

func (s *jwtService) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	claims := &model.JWTPayloadClaim{
		UserId: userID.String(),
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.AccessTokenLifetime)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: s.cfg.AppName,
			Subject: userID.String(),
		},
	}

	token := jwt.NewWithClaims(s.cfg.JwtSigningMethod, claims)
	return token.SignedString(s.cfg.JwtSignatureKey)
}

func (s *jwtService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.userRepo.FindById(ctx, userID)
	if err != nil || user == nil {
		return "", errors.New("failed to get user for refresh token")
	}

	claims := &model.JWTPayloadClaim{
		UserId: userID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.AppName,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(s.cfg.JwtSigningMethod, claims)
	tokenStr, err := token.SignedString(s.cfg.JwtSignatureKey)
	if err != nil {
		return "", err
	}

	rt := &model.RefreshToken{
		ID:              uuid.New(),
		UserID:          userID,
		Token:           tokenStr,
		Expires_At:      time.Now().Add(7 * 24 * time.Hour),
		Revoked:         false,
		Created_At:      time.Now(),
		Updated_At:      time.Now(),
	}

	if err := s.refreshToken.Save(ctx, rt); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *jwtService) ValidateRefreshToken(ctx context.Context, tokenStr string) (*model.JWTPayloadClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != s.cfg.JwtSigningMethod {
			return nil, errors.New("unexpected signing method")
		}
		return s.cfg.JwtSignatureKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(*model.JWTPayloadClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	rt, err := s.refreshToken.FindByToken(ctx, tokenStr)
	if err != nil || rt == nil {
		return nil, errors.New("refresh token not found")
	}
	if rt.Revoked || rt.Expires_At.Before(time.Now()) {
		return nil, errors.New("refresh token has expired or been revoked")
	}

	return claims, nil
}


func (s *jwtService) ValidateAccessToken(tokenStr string) (*model.JWTPayloadClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.JWTPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != s.cfg.JwtSigningMethod{
			return nil, errors.New("Unexpected signing method")
		}
		return s.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*model.JWTPayloadClaim)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token claims")
	}

	return claims, nil
}

func (s *jwtService) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	return s.refreshToken.RevokeById(ctx, tokenID, time.Now())
}

func (s *jwtService) RevokeRefreshTokenByToken(ctx context.Context, token string) error {
	return s.refreshToken.RevokeByToken(ctx, token, time.Now())
}


func (s *jwtService) DeleteExpiredToken(ctx context.Context) error {
	return  s.refreshToken.DeleteExpired(ctx)
}