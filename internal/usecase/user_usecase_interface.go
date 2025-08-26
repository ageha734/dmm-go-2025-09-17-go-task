package usecase

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, req CreateUserRequest, ipAddress, userAgent string) (*entity.User, error)
	GetUser(ctx context.Context, userID uint) (*entity.User, error)
	GetUsers(ctx context.Context, page, limit int) (*UserListResponse, error)
	UpdateUser(ctx context.Context, userID uint, req UpdateUserRequest, requestUserID uint, ipAddress, userAgent string) (*entity.User, error)
	DeleteUser(ctx context.Context, userID uint, requestUserID uint, ipAddress, userAgent string) error
	GetUserStats(ctx context.Context) (map[string]interface{}, error)
	HealthCheck(ctx context.Context) map[string]string
}
