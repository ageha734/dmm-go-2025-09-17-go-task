package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
)

type FraudDomainService struct {
	securityEventRepo     repository.SecurityEventRepository
	ipBlacklistRepo       repository.IPBlacklistRepository
	loginAttemptRepo      repository.LoginAttemptRepository
	rateLimitRuleRepo     repository.RateLimitRuleRepository
	userSessionRepo       repository.UserSessionRepository
	deviceFingerprintRepo repository.DeviceFingerprintRepository
}

func NewFraudDomainService(
	securityEventRepo repository.SecurityEventRepository,
	ipBlacklistRepo repository.IPBlacklistRepository,
	loginAttemptRepo repository.LoginAttemptRepository,
	rateLimitRuleRepo repository.RateLimitRuleRepository,
	userSessionRepo repository.UserSessionRepository,
	deviceFingerprintRepo repository.DeviceFingerprintRepository,
) *FraudDomainService {
	return &FraudDomainService{
		securityEventRepo:     securityEventRepo,
		ipBlacklistRepo:       ipBlacklistRepo,
		loginAttemptRepo:      loginAttemptRepo,
		rateLimitRuleRepo:     rateLimitRuleRepo,
		userSessionRepo:       userSessionRepo,
		deviceFingerprintRepo: deviceFingerprintRepo,
	}
}

func (s *FraudDomainService) AnalyzeFraud(ctx context.Context, userID *uint, email, ipAddress, userAgent string) (*entity.FraudAnalysis, error) {
	var riskScore float64
	var factors []string

	isBlacklisted, err := s.ipBlacklistRepo.IsBlacklisted(ctx, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP blacklist: %w", err)
	}
	if isBlacklisted {
		riskScore += 0.8
		factors = append(factors, "IP address is blacklisted")
	}

	since := time.Now().Add(-15 * time.Minute)
	failedAttempts, err := s.loginAttemptRepo.CountFailedAttempts(ctx, email, since)
	if err != nil {
		return nil, fmt.Errorf("failed to count failed attempts: %w", err)
	}
	if failedAttempts >= 5 {
		riskScore += 0.6
		factors = append(factors, fmt.Sprintf("Multiple failed login attempts: %d", failedAttempts))
	} else if failedAttempts >= 3 {
		riskScore += 0.3
		factors = append(factors, fmt.Sprintf("Some failed login attempts: %d", failedAttempts))
	}

	if userID != nil {

		fingerprint := fmt.Sprintf("%s_%s", ipAddress, userAgent)
		isTrusted, err := s.deviceFingerprintRepo.IsTrustedDevice(ctx, *userID, fingerprint)
		if err != nil {
			return nil, fmt.Errorf("failed to check device trust: %w", err)
		}
		if !isTrusted {
			riskScore += 0.2
			factors = append(factors, "Unknown device")
		}
	}

	ipAttempts, err := s.loginAttemptRepo.GetByIP(ctx, ipAddress, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP attempts: %w", err)
	}
	if len(ipAttempts) >= 10 {
		riskScore += 0.4
		factors = append(factors, fmt.Sprintf("High frequency requests from IP: %d", len(ipAttempts)))
	}

	if riskScore > 1.0 {
		riskScore = 1.0
	}

	return entity.NewFraudAnalysis(riskScore, factors), nil
}

func (s *FraudDomainService) RecordLoginAttempt(ctx context.Context, email, ipAddress, userAgent string, success bool, failReason string) error {
	attempt := entity.NewLoginAttempt(email, ipAddress, userAgent, success, failReason)
	return s.loginAttemptRepo.Create(ctx, attempt)
}

func (s *FraudDomainService) CreateSecurityEvent(ctx context.Context, userID *uint, eventType, description, ipAddress, userAgent, severity string) error {
	event := entity.NewSecurityEvent(userID, eventType, description, ipAddress, userAgent, severity)
	return s.securityEventRepo.Create(ctx, event)
}

func (s *FraudDomainService) AddIPToBlacklist(ctx context.Context, ipAddress, reason string, expiresAt *time.Time) error {
	blacklist := entity.NewIPBlacklist(ipAddress, reason, expiresAt)
	return s.ipBlacklistRepo.Create(ctx, blacklist)
}

func (s *FraudDomainService) RemoveIPFromBlacklist(ctx context.Context, ipAddress string) error {
	blacklist, err := s.ipBlacklistRepo.GetByIP(ctx, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to get IP blacklist: %w", err)
	}

	blacklist.Deactivate()
	return s.ipBlacklistRepo.Update(ctx, blacklist)
}

func (s *FraudDomainService) CheckRateLimit(ctx context.Context, resource string) (*entity.RateLimitRule, error) {
	return s.rateLimitRuleRepo.GetByResource(ctx, resource)
}

func (s *FraudDomainService) CreateRateLimitRule(ctx context.Context, name, resource string, maxRequests, windowSize int) error {
	rule := entity.NewRateLimitRule(name, resource, maxRequests, windowSize)
	return s.rateLimitRuleRepo.Create(ctx, rule)
}

func (s *FraudDomainService) UpdateRateLimitRule(ctx context.Context, ruleID uint, maxRequests, windowSize int) error {
	rule, err := s.rateLimitRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("failed to get rate limit rule: %w", err)
	}

	rule.UpdateRule(maxRequests, windowSize)
	return s.rateLimitRuleRepo.Update(ctx, rule)
}

func (s *FraudDomainService) DeleteRateLimitRule(ctx context.Context, ruleID uint) error {
	rule, err := s.rateLimitRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("failed to get rate limit rule: %w", err)
	}

	rule.Deactivate()
	return s.rateLimitRuleRepo.Update(ctx, rule)
}

func (s *FraudDomainService) CreateUserSession(ctx context.Context, userID uint, sessionID, ipAddress, userAgent string, expiresAt time.Time) error {
	session := entity.NewUserSession(userID, sessionID, ipAddress, userAgent, expiresAt)
	return s.userSessionRepo.Create(ctx, session)
}

func (s *FraudDomainService) ValidateUserSession(ctx context.Context, sessionID string) (*entity.UserSession, error) {
	session, err := s.userSessionRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	if !session.IsValid() {
		return nil, fmt.Errorf("session is invalid or expired")
	}

	return session, nil
}

func (s *FraudDomainService) DeactivateUserSession(ctx context.Context, sessionID string) error {
	session, err := s.userSessionRepo.GetBySessionID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get user session: %w", err)
	}

	session.Deactivate()
	return s.userSessionRepo.Update(ctx, session)
}

func (s *FraudDomainService) RecordDeviceFingerprint(ctx context.Context, userID uint, fingerprint string, deviceInfo *string) error {
	existing, err := s.deviceFingerprintRepo.GetByFingerprint(ctx, fingerprint)
	if err == nil {
		existing.UpdateLastSeen()
		return s.deviceFingerprintRepo.Update(ctx, existing)
	}

	deviceFingerprint := entity.NewDeviceFingerprint(userID, fingerprint, deviceInfo)
	return s.deviceFingerprintRepo.Create(ctx, deviceFingerprint)
}

func (s *FraudDomainService) TrustDevice(ctx context.Context, userID uint, fingerprint string) error {
	deviceFingerprint, err := s.deviceFingerprintRepo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to get device fingerprint: %w", err)
	}

	if deviceFingerprint.UserID != userID {
		return fmt.Errorf("device does not belong to user")
	}

	deviceFingerprint.Trust()
	return s.deviceFingerprintRepo.Update(ctx, deviceFingerprint)
}

func (s *FraudDomainService) CleanupExpiredData(ctx context.Context) error {
	if err := s.ipBlacklistRepo.CleanupExpired(ctx); err != nil {
		return fmt.Errorf("failed to cleanup expired IP blacklist: %w", err)
	}

	if err := s.userSessionRepo.CleanupExpired(ctx); err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	before := time.Now().AddDate(0, 0, -30)
	if err := s.loginAttemptRepo.CleanupOld(ctx, before); err != nil {
		return fmt.Errorf("failed to cleanup old login attempts: %w", err)
	}

	return nil
}
