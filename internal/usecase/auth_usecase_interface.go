package usecase

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
)

type AuthUsecaseInterface interface {
	Register(ctx context.Context, req RegisterRequest, ipAddress, userAgent string) (*LoginResponse, error)
	Login(ctx context.Context, req LoginRequest, ipAddress, userAgent string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (*LoginResponse, error)
	ChangePassword(ctx context.Context, userID uint, req ChangePasswordRequest, ipAddress, userAgent string) error
	Logout(ctx context.Context, userID uint, sessionID, ipAddress, userAgent string) error
	ValidateToken(tokenString string) (*service.JWTClaims, error)
}
