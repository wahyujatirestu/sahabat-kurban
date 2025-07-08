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
)

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error)
	ValidateAccessToken(tokenStr string) (*model.JWTPayloadClaim, error)
	RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpiredToken(ctx context.Context) error
}

type jwtService struct {
	cfg		*config.Config
	refreshToken repository.RefreshTokenRepository
}

func NewJWTServie(cfg *config.Config, repo repository.RefreshTokenRepository, ) JWTService {
	return &jwtService {
		cfg: cfg, refreshToken: repo,
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

func (s *jwtService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (*model.RefreshToken, error) {
	tokenUUID := uuid.New()
	tokenStr := tokenUUID.String()

	rt := &model.RefreshToken{
		ID: uuid.New(),
		UserID: userID,
		Token: tokenStr,
		Expires_At: time.Now().Add(7 * 24 * time.Hour),
		Revoked: false,
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	if err := s.refreshToken.Save(ctx, rt); err != nil {
		return nil, err
	}

	return rt, nil
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

func (s *jwtService) DeleteExpiredToken(ctx context.Context) error {
	return  s.refreshToken.DeleteExpired(ctx)
}