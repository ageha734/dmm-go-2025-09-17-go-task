package service

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

type FraudDomainServiceInterface interface {
	AnalyzeFraud(ctx context.Context, userID *uint, email, ipAddress, userAgent string) (*entity.FraudAnalysis, error)
	CreateSecurityEvent(ctx context.Context, userID *uint, eventType, description, ipAddress, userAgent, severity string) error
	RecordLoginAttempt(ctx context.Context, email, ipAddress, userAgent string, success bool, failureReason string) error
	DeactivateUserSession(ctx context.Context, sessionID string) error
	GetFraudStats(ctx context.Context) (map[string]interface{}, error)

	AddIPToBlacklist(ctx context.Context, ip, reason, clientIP, userAgent string) error
	RemoveIPFromBlacklist(ctx context.Context, ip, clientIP, userAgent string) error
	GetBlacklistedIPs(ctx context.Context, page, limit int) (interface{}, error)

	GetSecurityEvents(ctx context.Context, page, limit int) (interface{}, error)

	CreateRateLimitRule(ctx context.Context, name, pattern string, maxRequests, windowSize int64) error
	UpdateRateLimitRule(ctx context.Context, id uint, name, pattern string, maxRequests, windowSize int64) error
	DeleteRateLimitRule(ctx context.Context, id uint) error
	GetRateLimitRules(ctx context.Context) (interface{}, error)

	GetActiveSessions(ctx context.Context) (interface{}, error)
	DeactivateSession(ctx context.Context, sessionID string) error

	GetDevices(ctx context.Context) (interface{}, error)
	TrustDevice(ctx context.Context, fingerprint string) error

	CleanupExpiredData(ctx context.Context) error
}
