package repository

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type AuthRepository interface {
	Create(ctx context.Context, auth *entity.Auth) error

	GetByUserID(ctx context.Context, userID uint) (*entity.Auth, error)

	GetByEmail(ctx context.Context, email string) (*entity.Auth, error)

	Update(ctx context.Context, auth *entity.Auth) error

	Delete(ctx context.Context, userID uint) error

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error

	GetByID(ctx context.Context, id uint) (*entity.Role, error)

	GetByName(ctx context.Context, name string) (*entity.Role, error)

	List(ctx context.Context) ([]*entity.Role, error)

	Update(ctx context.Context, role *entity.Role) error

	Delete(ctx context.Context, id uint) error

	AssignToUser(ctx context.Context, userID, roleID uint) error

	RemoveFromUser(ctx context.Context, userID, roleID uint) error

	GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error

	GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error)

	Update(ctx context.Context, token *entity.RefreshToken) error

	RevokeByUserID(ctx context.Context, userID uint) error

	DeleteExpired(ctx context.Context) error
}
