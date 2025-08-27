package usecase

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type FraudUsecaseInterface interface {
	AddIPToBlacklist(ctx context.Context, ip string, reason string, adminID uint) error
	RemoveIPFromBlacklist(ctx context.Context, ip string) error
	GetBlacklistedIPs(ctx context.Context) ([]*entity.IPBlacklist, error)
	GetSecurityEvents(ctx context.Context, limit, offset int) ([]*entity.SecurityEvent, error)
	CreateSecurityEvent(ctx context.Context, req *dto.CreateSecurityEventRequest) error
	CreateRateLimitRule(ctx context.Context, req *dto.CreateRateLimitRuleRequest) error
	UpdateRateLimitRule(ctx context.Context, id uint, req *dto.UpdateRateLimitRuleRequest) error
	DeleteRateLimitRule(ctx context.Context, id uint) error
	GetRateLimitRules(ctx context.Context) ([]*entity.RateLimitRule, error)
	GetActiveSessions(ctx context.Context) ([]*entity.UserSession, error)
	DeactivateSession(ctx context.Context, sessionID string) error
	GetDevices(ctx context.Context) ([]*entity.DeviceFingerprint, error)
	TrustDevice(ctx context.Context, fingerprint string) error
	CleanupExpiredData(ctx context.Context) error
}
