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
}
