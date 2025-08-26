package repository

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type MembershipTierRepository interface {
	Create(ctx context.Context, tier *entity.MembershipTier) error

	GetByID(ctx context.Context, id uint) (*entity.MembershipTier, error)

	GetByName(ctx context.Context, name string) (*entity.MembershipTier, error)

	List(ctx context.Context) ([]*entity.MembershipTier, error)

	Update(ctx context.Context, tier *entity.MembershipTier) error

	Delete(ctx context.Context, id uint) error
}

type UserMembershipRepository interface {
	Create(ctx context.Context, membership *entity.UserMembership) error

	GetByUserID(ctx context.Context, userID uint) (*entity.UserMembership, error)

	Update(ctx context.Context, membership *entity.UserMembership) error

	Delete(ctx context.Context, userID uint) error

	GetStats(ctx context.Context) (map[string]interface{}, error)

	List(ctx context.Context, offset, limit int) ([]*entity.UserMembership, int64, error)
}

type PointTransactionRepository interface {
	Create(ctx context.Context, transaction *entity.PointTransaction) error

	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.PointTransaction, int64, error)

	GetByID(ctx context.Context, id uint) (*entity.PointTransaction, error)

	List(ctx context.Context, offset, limit int) ([]*entity.PointTransaction, int64, error)

	ExpirePoints(ctx context.Context) error
}

type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error

	GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error)

	Update(ctx context.Context, profile *entity.UserProfile) error

	Delete(ctx context.Context, userID uint) error
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error

	GetByUserID(ctx context.Context, userID uint, offset, limit int, unreadOnly bool) ([]*entity.Notification, int64, error)

	GetByID(ctx context.Context, id uint) (*entity.Notification, error)

	Update(ctx context.Context, notification *entity.Notification) error

	Delete(ctx context.Context, id uint) error

	MarkAsRead(ctx context.Context, userID, notificationID uint) error

	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
}
