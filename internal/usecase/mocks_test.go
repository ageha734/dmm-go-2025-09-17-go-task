package usecase_test

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/stretchr/testify/mock"
)

type MockAuthDomainService struct {
	mock.Mock
}

func (m *MockAuthDomainService) Register(ctx context.Context, name, email, password string, age int) (*entity.User, error) {
	args := m.Called(ctx, name, email, password, age)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthDomainService) Login(ctx context.Context, email, password string) (*entity.Auth, []string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var auth *entity.Auth
	if a, ok := args.Get(0).(*entity.Auth); ok {
		auth = a
	}
	var roles []string
	if r, ok := args.Get(1).([]string); ok {
		roles = r
	}
	return auth, roles, args.Error(2)
}

func (m *MockAuthDomainService) GenerateAccessToken(userID uint, email string, roles []string) (string, error) {
	args := m.Called(userID, email, roles)
	return args.String(0), args.Error(1)
}

func (m *MockAuthDomainService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthDomainService) ValidateToken(tokenString string) (*service.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if claims, ok := args.Get(0).(*service.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthDomainService) RefreshToken(ctx context.Context, refreshTokenStr string) (*entity.Auth, []string, string, error) {
	args := m.Called(ctx, refreshTokenStr)
	if args.Get(0) == nil {
		return nil, nil, "", args.Error(3)
	}
	var auth *entity.Auth
	if a, ok := args.Get(0).(*entity.Auth); ok {
		auth = a
	}
	var roles []string
	if r, ok := args.Get(1).([]string); ok {
		roles = r
	}
	return auth, roles, args.String(2), args.Error(3)
}

func (m *MockAuthDomainService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthDomainService) Logout(ctx context.Context, userID uint, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

type MockFraudDomainService struct {
	mock.Mock
}

func (m *MockFraudDomainService) AnalyzeFraud(ctx context.Context, userID *uint, email, ipAddress, userAgent string) (*entity.FraudAnalysis, error) {
	args := m.Called(ctx, userID, email, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if analysis, ok := args.Get(0).(*entity.FraudAnalysis); ok {
		return analysis, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFraudDomainService) CreateSecurityEvent(ctx context.Context, userID *uint, eventType, description, ipAddress, userAgent, severity string) error {
	args := m.Called(ctx, userID, eventType, description, ipAddress, userAgent, severity)
	return args.Error(0)
}

func (m *MockFraudDomainService) RecordLoginAttempt(ctx context.Context, email, ipAddress, userAgent string, success bool, failureReason string) error {
	args := m.Called(ctx, email, ipAddress, userAgent, success, failureReason)
	return args.Error(0)
}

func (m *MockFraudDomainService) DeactivateUserSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetFraudStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if stats, ok := args.Get(0).(map[string]interface{}); ok {
		return stats, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFraudDomainService) AddIPToBlacklist(ctx context.Context, ip, reason, clientIP, userAgent string) error {
	args := m.Called(ctx, ip, reason, clientIP, userAgent)
	return args.Error(0)
}

func (m *MockFraudDomainService) RemoveIPFromBlacklist(ctx context.Context, ip, clientIP, userAgent string) error {
	args := m.Called(ctx, ip, clientIP, userAgent)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetBlacklistedIPs(ctx context.Context, page, limit int) (interface{}, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) GetSecurityEvents(ctx context.Context, page, limit int) (interface{}, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) CreateRateLimitRule(ctx context.Context, name, pattern string, maxRequests, windowSize int64) error {
	args := m.Called(ctx, name, pattern, maxRequests, windowSize)
	return args.Error(0)
}

func (m *MockFraudDomainService) UpdateRateLimitRule(ctx context.Context, id uint, name, pattern string, maxRequests, windowSize int64) error {
	args := m.Called(ctx, id, name, pattern, maxRequests, windowSize)
	return args.Error(0)
}

func (m *MockFraudDomainService) DeleteRateLimitRule(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetRateLimitRules(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) GetActiveSessions(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) DeactivateSession(ctx context.Context, sessionID string) error {
	return m.DeactivateUserSession(ctx, sessionID)
}

func (m *MockFraudDomainService) GetDevices(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

func (m *MockFraudDomainService) TrustDevice(ctx context.Context, fingerprint string) error {
	args := m.Called(ctx, fingerprint)
	return args.Error(0)
}

func (m *MockFraudDomainService) CleanupExpiredData(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
