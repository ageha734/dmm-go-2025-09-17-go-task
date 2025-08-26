package repository

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error

	GetByID(ctx context.Context, id uint) (*entity.User, error)

	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	Update(ctx context.Context, user *entity.User) error

	Delete(ctx context.Context, id uint) error

	List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
