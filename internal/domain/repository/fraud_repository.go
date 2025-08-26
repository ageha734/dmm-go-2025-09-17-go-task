package repository

import (
	"context"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type SecurityEventRepository interface {
	Create(ctx context.Context, event *entity.SecurityEvent) error

	GetByID(ctx context.Context, id uint) (*entity.SecurityEvent, error)

	List(ctx context.Context, offset, limit int) ([]*entity.SecurityEvent, int64, error)

	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.SecurityEvent, int64, error)

	GetBySeverity(ctx context.Context, severity string, offset, limit int) ([]*entity.SecurityEvent, int64, error)

	Delete(ctx context.Context, id uint) error
}

type IPBlacklistRepository interface {
	Create(ctx context.Context, blacklist *entity.IPBlacklist) error

	GetByIP(ctx context.Context, ipAddress string) (*entity.IPBlacklist, error)

	List(ctx context.Context, offset, limit int) ([]*entity.IPBlacklist, int64, error)

	Update(ctx context.Context, blacklist *entity.IPBlacklist) error

	Delete(ctx context.Context, ipAddress string) error

	IsBlacklisted(ctx context.Context, ipAddress string) (bool, error)

	CleanupExpired(ctx context.Context) error
}

type LoginAttemptRepository interface {
	Create(ctx context.Context, attempt *entity.LoginAttempt) error

	GetByEmail(ctx context.Context, email string, since time.Time) ([]*entity.LoginAttempt, error)

	GetByIP(ctx context.Context, ipAddress string, since time.Time) ([]*entity.LoginAttempt, error)

	CountFailedAttempts(ctx context.Context, email string, since time.Time) (int64, error)

	List(ctx context.Context, offset, limit int) ([]*entity.LoginAttempt, int64, error)

	Delete(ctx context.Context, id uint) error

	CleanupOld(ctx context.Context, before time.Time) error
}

type RateLimitRuleRepository interface {
	Create(ctx context.Context, rule *entity.RateLimitRule) error

	GetByID(ctx context.Context, id uint) (*entity.RateLimitRule, error)

	GetByResource(ctx context.Context, resource string) (*entity.RateLimitRule, error)

	List(ctx context.Context) ([]*entity.RateLimitRule, error)

	Update(ctx context.Context, rule *entity.RateLimitRule) error

	Delete(ctx context.Context, id uint) error

	GetActiveRules(ctx context.Context) ([]*entity.RateLimitRule, error)
}

type UserSessionRepository interface {
	Create(ctx context.Context, session *entity.UserSession) error

	GetBySessionID(ctx context.Context, sessionID string) (*entity.UserSession, error)

	GetByUserID(ctx context.Context, userID uint) ([]*entity.UserSession, error)

	Update(ctx context.Context, session *entity.UserSession) error

	Delete(ctx context.Context, sessionID string) error

	DeactivateByUserID(ctx context.Context, userID uint) error

	CleanupExpired(ctx context.Context) error
}

type DeviceFingerprintRepository interface {
	Create(ctx context.Context, fingerprint *entity.DeviceFingerprint) error

	GetByFingerprint(ctx context.Context, fingerprint string) (*entity.DeviceFingerprint, error)

	GetByUserID(ctx context.Context, userID uint) ([]*entity.DeviceFingerprint, error)

	Update(ctx context.Context, fingerprint *entity.DeviceFingerprint) error

	Delete(ctx context.Context, id uint) error

	IsTrustedDevice(ctx context.Context, userID uint, fingerprint string) (bool, error)
}
