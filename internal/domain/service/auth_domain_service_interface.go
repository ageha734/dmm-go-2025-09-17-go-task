package service

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type AuthDomainServiceInterface interface {
	Register(ctx context.Context, name, email, password string, age int) (*entity.User, error)
	Login(ctx context.Context, email, password string) (*entity.Auth, []string, error)
	GenerateAccessToken(userID uint, email string, roles []string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID uint) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(ctx context.Context, refreshTokenStr string) (*entity.Auth, []string, string, error)
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	Logout(ctx context.Context, userID uint, token string) error
}
