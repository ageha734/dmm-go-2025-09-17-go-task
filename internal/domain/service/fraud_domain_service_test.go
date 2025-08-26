package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecurityEventRepository struct {
	mock.Mock
}

func (m *MockSecurityEventRepository) Create(ctx context.Context, event *entity.SecurityEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockSecurityEventRepository) GetByID(ctx context.Context, id uint) (*entity.SecurityEvent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if event, ok := args.Get(0).(*entity.SecurityEvent); ok {
		return event, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSecurityEventRepository) List(ctx context.Context, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	args := m.Called(ctx, offset, limit)
	var events []*entity.SecurityEvent
	if e, ok := args.Get(0).([]*entity.SecurityEvent); ok {
		events = e
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return events, total, args.Error(2)
}

func (m *MockSecurityEventRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	var events []*entity.SecurityEvent
	if e, ok := args.Get(0).([]*entity.SecurityEvent); ok {
		events = e
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return events, total, args.Error(2)
}

func (m *MockSecurityEventRepository) GetBySeverity(ctx context.Context, severity string, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	args := m.Called(ctx, severity, offset, limit)
	var events []*entity.SecurityEvent
	if e, ok := args.Get(0).([]*entity.SecurityEvent); ok {
		events = e
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return events, total, args.Error(2)
}

func (m *MockSecurityEventRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockIPBlacklistRepository struct {
	mock.Mock
}

func (m *MockIPBlacklistRepository) Create(ctx context.Context, blacklist *entity.IPBlacklist) error {
	args := m.Called(ctx, blacklist)
	return args.Error(0)
}

func (m *MockIPBlacklistRepository) GetByID(ctx context.Context, id uint) (*entity.IPBlacklist, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if blacklist, ok := args.Get(0).(*entity.IPBlacklist); ok {
		return blacklist, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockIPBlacklistRepository) GetByIP(ctx context.Context, ipAddress string) (*entity.IPBlacklist, error) {
	args := m.Called(ctx, ipAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if blacklist, ok := args.Get(0).(*entity.IPBlacklist); ok {
		return blacklist, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockIPBlacklistRepository) IsBlacklisted(ctx context.Context, ipAddress string) (bool, error) {
	args := m.Called(ctx, ipAddress)
	return args.Bool(0), args.Error(1)
}

func (m *MockIPBlacklistRepository) Update(ctx context.Context, blacklist *entity.IPBlacklist) error {
	if blacklist == nil {
		return errors.New("blacklist is nil")
	}
	args := m.Called(ctx, blacklist)
	return args.Error(0)
}

func (m *MockIPBlacklistRepository) List(ctx context.Context, offset, limit int) ([]*entity.IPBlacklist, int64, error) {
	args := m.Called(ctx, offset, limit)
	var blacklists []*entity.IPBlacklist
	if b, ok := args.Get(0).([]*entity.IPBlacklist); ok {
		blacklists = b
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return blacklists, total, args.Error(2)
}

func (m *MockIPBlacklistRepository) Delete(ctx context.Context, ipAddress string) error {
	args := m.Called(ctx, ipAddress)
	return args.Error(0)
}

func (m *MockIPBlacklistRepository) CleanupExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockLoginAttemptRepository struct {
	mock.Mock
}

func (m *MockLoginAttemptRepository) Create(ctx context.Context, attempt *entity.LoginAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockLoginAttemptRepository) GetByID(ctx context.Context, id uint) (*entity.LoginAttempt, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if attempt, ok := args.Get(0).(*entity.LoginAttempt); ok {
		return attempt, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoginAttemptRepository) GetByEmail(ctx context.Context, email string, since time.Time) ([]*entity.LoginAttempt, error) {
	args := m.Called(ctx, email, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if attempts, ok := args.Get(0).([]*entity.LoginAttempt); ok {
		return attempts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoginAttemptRepository) GetByIP(ctx context.Context, ipAddress string, since time.Time) ([]*entity.LoginAttempt, error) {
	args := m.Called(ctx, ipAddress, since)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if attempts, ok := args.Get(0).([]*entity.LoginAttempt); ok {
		return attempts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLoginAttemptRepository) CountFailedAttempts(ctx context.Context, email string, since time.Time) (int64, error) {
	args := m.Called(ctx, email, since)
	var count int64
	if c, ok := args.Get(0).(int64); ok {
		count = c
	}
	return count, args.Error(1)
}

func (m *MockLoginAttemptRepository) List(ctx context.Context, offset, limit int) ([]*entity.LoginAttempt, int64, error) {
	args := m.Called(ctx, offset, limit)
	var attempts []*entity.LoginAttempt
	if a, ok := args.Get(0).([]*entity.LoginAttempt); ok {
		attempts = a
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return attempts, total, args.Error(2)
}

func (m *MockLoginAttemptRepository) Update(ctx context.Context, attempt *entity.LoginAttempt) error {
	if attempt == nil {
		return errors.New("attempt is nil")
	}
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockLoginAttemptRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLoginAttemptRepository) CleanupOld(ctx context.Context, before time.Time) error {
	args := m.Called(ctx, before)
	return args.Error(0)
}

type MockRateLimitRuleRepository struct {
	mock.Mock
}

func (m *MockRateLimitRuleRepository) Create(ctx context.Context, rule *entity.RateLimitRule) error {
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRateLimitRuleRepository) GetByID(ctx context.Context, id uint) (*entity.RateLimitRule, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if rule, ok := args.Get(0).(*entity.RateLimitRule); ok {
		return rule, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRateLimitRuleRepository) GetByResource(ctx context.Context, resource string) (*entity.RateLimitRule, error) {
	args := m.Called(ctx, resource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if rule, ok := args.Get(0).(*entity.RateLimitRule); ok {
		return rule, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRateLimitRuleRepository) Update(ctx context.Context, rule *entity.RateLimitRule) error {
	if rule == nil {
		return errors.New("rule is nil")
	}
	args := m.Called(ctx, rule)
	return args.Error(0)
}

func (m *MockRateLimitRuleRepository) List(ctx context.Context) ([]*entity.RateLimitRule, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if rules, ok := args.Get(0).([]*entity.RateLimitRule); ok {
		return rules, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRateLimitRuleRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRateLimitRuleRepository) GetActiveRules(ctx context.Context) ([]*entity.RateLimitRule, error) {
	args := m.Called(ctx)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	if rules, ok := result.([]*entity.RateLimitRule); ok {
		return rules, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUserSessionRepository struct {
	mock.Mock
}

func (m *MockUserSessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserSessionRepository) GetByID(ctx context.Context, id uint) (*entity.UserSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if session, ok := args.Get(0).(*entity.UserSession); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserSessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*entity.UserSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if session, ok := args.Get(0).(*entity.UserSession); ok {
		return session, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserSessionRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.UserSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if sessions, ok := args.Get(0).([]*entity.UserSession); ok {
		return sessions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserSessionRepository) Update(ctx context.Context, session *entity.UserSession) error {
	if session == nil {
		return errors.New("session is nil")
	}
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockUserSessionRepository) Delete(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockUserSessionRepository) DeactivateByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserSessionRepository) CleanupExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockDeviceFingerprintRepository struct {
	mock.Mock
}

func (m *MockDeviceFingerprintRepository) Create(ctx context.Context, fingerprint *entity.DeviceFingerprint) error {
	args := m.Called(ctx, fingerprint)
	return args.Error(0)
}

func (m *MockDeviceFingerprintRepository) GetByID(ctx context.Context, id uint) (*entity.DeviceFingerprint, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if fingerprint, ok := args.Get(0).(*entity.DeviceFingerprint); ok {
		return fingerprint, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceFingerprintRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*entity.DeviceFingerprint, error) {
	args := m.Called(ctx, fingerprint)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if fp, ok := args.Get(0).(*entity.DeviceFingerprint); ok {
		return fp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceFingerprintRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.DeviceFingerprint, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if fingerprints, ok := args.Get(0).([]*entity.DeviceFingerprint); ok {
		return fingerprints, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceFingerprintRepository) IsTrustedDevice(ctx context.Context, userID uint, fingerprint string) (bool, error) {
	args := m.Called(ctx, userID, fingerprint)
	return args.Bool(0), args.Error(1)
}

func (m *MockDeviceFingerprintRepository) Update(ctx context.Context, fingerprint *entity.DeviceFingerprint) error {
	if fingerprint == nil {
		return errors.New("fingerprint is nil")
	}
	args := m.Called(ctx, fingerprint)
	return args.Error(0)
}

func (m *MockDeviceFingerprintRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupFraudDomainService() (*service.FraudDomainService, *MockSecurityEventRepository, *MockIPBlacklistRepository, *MockLoginAttemptRepository, *MockRateLimitRuleRepository, *MockUserSessionRepository, *MockDeviceFingerprintRepository) {
	mockSecurityEventRepo := &MockSecurityEventRepository{}
	mockIPBlacklistRepo := &MockIPBlacklistRepository{}
	mockLoginAttemptRepo := &MockLoginAttemptRepository{}
	mockRateLimitRuleRepo := &MockRateLimitRuleRepository{}
	mockUserSessionRepo := &MockUserSessionRepository{}
	mockDeviceFingerprintRepo := &MockDeviceFingerprintRepository{}

	service := service.NewFraudDomainService(
		mockSecurityEventRepo,
		mockIPBlacklistRepo,
		mockLoginAttemptRepo,
		mockRateLimitRuleRepo,
		mockUserSessionRepo,
		mockDeviceFingerprintRepo,
	)

	return service, mockSecurityEventRepo, mockIPBlacklistRepo, mockLoginAttemptRepo, mockRateLimitRuleRepo, mockUserSessionRepo, mockDeviceFingerprintRepo
}

func TestNewFraudDomainService(t *testing.T) {
	t.Run("新しいFraudDomainServiceを正常に作成できる", func(t *testing.T) {
		service, _, _, _, _, _, _ := setupFraudDomainService()

		assert.NotNil(t, service)
	})
}
