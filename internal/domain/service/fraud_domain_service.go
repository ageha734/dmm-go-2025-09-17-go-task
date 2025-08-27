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

func (s *FraudDomainService) AddIPToBlacklist(ctx context.Context, ip, reason, clientIP, userAgent string) error {
	blacklist := entity.NewIPBlacklist(ip, reason, nil)
	err := s.ipBlacklistRepo.Create(ctx, blacklist)
	if err != nil {
		return err
	}

	return s.CreateSecurityEvent(ctx, nil, "IP_BLACKLISTED", fmt.Sprintf("IP %s blacklisted: %s", ip, reason), clientIP, userAgent, "MEDIUM")
}

func (s *FraudDomainService) RemoveIPFromBlacklist(ctx context.Context, ip, clientIP, userAgent string) error {
	blacklist, err := s.ipBlacklistRepo.GetByIP(ctx, ip)
	if err != nil {
		return fmt.Errorf("failed to get IP blacklist: %w", err)
	}

	blacklist.Deactivate()
	err = s.ipBlacklistRepo.Update(ctx, blacklist)
	if err != nil {
		return err
	}

	return s.CreateSecurityEvent(ctx, nil, "IP_UNBLACKLISTED", fmt.Sprintf("IP %s removed from blacklist", ip), clientIP, userAgent, "LOW")
}

func (s *FraudDomainService) GetBlacklistedIPs(ctx context.Context, page, limit int) (interface{}, error) {
	ips, total, err := s.ipBlacklistRepo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklisted IPs: %w", err)
	}

	return map[string]interface{}{
		"ips":   ips,
		"total": total,
		"page":  page,
		"limit": limit,
	}, nil
}

func (s *FraudDomainService) GetSecurityEvents(ctx context.Context, page, limit int) (interface{}, error) {
	events, total, err := s.securityEventRepo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	return map[string]interface{}{
		"events": events,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}, nil
}

func (s *FraudDomainService) CheckRateLimit(ctx context.Context, resource string) (*entity.RateLimitRule, error) {
	return s.rateLimitRuleRepo.GetByResource(ctx, resource)
}

func (s *FraudDomainService) CreateRateLimitRule(ctx context.Context, name, pattern string, maxRequests, windowSize int64) error {
	rule := entity.NewRateLimitRule(name, pattern, int(maxRequests), int(windowSize))
	return s.rateLimitRuleRepo.Create(ctx, rule)
}

func (s *FraudDomainService) UpdateRateLimitRule(ctx context.Context, ruleID uint, _, _ string, maxRequests, windowSize int64) error {
	rule, err := s.rateLimitRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return fmt.Errorf("failed to get rate limit rule: %w", err)
	}

	rule.UpdateRule(int(maxRequests), int(windowSize))
	return s.rateLimitRuleRepo.Update(ctx, rule)
}

func (s *FraudDomainService) GetRateLimitRules(ctx context.Context) (interface{}, error) {
	rules, err := s.rateLimitRuleRepo.GetActiveRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit rules: %w", err)
	}

	return map[string]interface{}{
		"rules": rules,
		"total": len(rules),
	}, nil
}

func (s *FraudDomainService) GetActiveSessions(ctx context.Context) (interface{}, error) {
	now := time.Now()

	sessionEvents, _, err := s.securityEventRepo.List(ctx, 1, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to get session events: %w", err)
	}

	activeSessions := []map[string]interface{}{}
	for _, event := range sessionEvents {
		if event.EventType == "SESSION_CREATED" || event.EventType == "LOGIN_SUCCESS" {
			sessionInfo := map[string]interface{}{
				"session_id": fmt.Sprintf("sess_%d", event.ID),
				"user_id":    event.UserID,
				"ip_address": event.IPAddress,
				"user_agent": event.UserAgent,
				"created_at": event.CreatedAt,
				"last_seen":  event.CreatedAt,
				"is_active":  true,
			}
			activeSessions = append(activeSessions, sessionInfo)
		}
	}

	return map[string]interface{}{
		"sessions":     activeSessions,
		"total":        len(activeSessions),
		"active_count": len(activeSessions),
		"timestamp":    now,
	}, nil
}

func (s *FraudDomainService) DeactivateSession(ctx context.Context, sessionID string) error {
	return s.DeactivateUserSession(ctx, sessionID)
}

func (s *FraudDomainService) GetDevices(ctx context.Context) (interface{}, error) {
	now := time.Now()
	since := now.Add(-30 * 24 * time.Hour)

	loginAttempts, err := s.getFilteredLoginAttempts(ctx, since)
	if err != nil {
		return nil, err
	}

	deviceMap := s.buildDeviceMap(loginAttempts)
	devices := s.convertDeviceMapToSlice(deviceMap)

	return map[string]interface{}{
		"devices":   devices,
		"total":     len(devices),
		"timestamp": now,
		"period":    "30 days",
	}, nil
}

func (s *FraudDomainService) getFilteredLoginAttempts(ctx context.Context, since time.Time) ([]*entity.LoginAttempt, error) {
	loginAttempts, _, err := s.loginAttemptRepo.List(ctx, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get login attempts: %w", err)
	}

	filteredAttempts := make([]*entity.LoginAttempt, 0, len(loginAttempts))
	for _, attempt := range loginAttempts {
		if attempt.CreatedAt.After(since) {
			filteredAttempts = append(filteredAttempts, attempt)
		}
	}
	return filteredAttempts, nil
}

func (s *FraudDomainService) buildDeviceMap(loginAttempts []*entity.LoginAttempt) map[string]map[string]interface{} {
	deviceMap := make(map[string]map[string]interface{})

	for _, attempt := range loginAttempts {
		fingerprint := fmt.Sprintf("%s_%s", attempt.IPAddress, attempt.UserAgent)

		if _, exists := deviceMap[fingerprint]; !exists {
			deviceMap[fingerprint] = s.createNewDeviceInfo(attempt, fingerprint)
		} else {
			s.updateExistingDeviceInfo(deviceMap[fingerprint], attempt)
		}
	}
	return deviceMap
}

func (s *FraudDomainService) createNewDeviceInfo(attempt *entity.LoginAttempt, fingerprint string) map[string]interface{} {
	deviceInfo := map[string]interface{}{
		"fingerprint":   fingerprint,
		"ip_address":    attempt.IPAddress,
		"user_agent":    attempt.UserAgent,
		"first_seen":    attempt.CreatedAt,
		"last_seen":     attempt.CreatedAt,
		"login_count":   1,
		"success_count": 0,
		"is_trusted":    false,
		"risk_level":    "unknown",
	}

	if attempt.Success {
		deviceInfo["success_count"] = 1
		deviceInfo["risk_level"] = "low"
	} else {
		deviceInfo["risk_level"] = "medium"
	}

	return deviceInfo
}

func (s *FraudDomainService) updateExistingDeviceInfo(device map[string]interface{}, attempt *entity.LoginAttempt) {
	device["last_seen"] = attempt.CreatedAt

	loginCount, ok := device["login_count"].(int)
	if !ok {
		loginCount = 0
	}
	device["login_count"] = loginCount + 1

	if attempt.Success {
		successCount, ok := device["success_count"].(int)
		if !ok {
			successCount = 0
		}
		device["success_count"] = successCount + 1
	}

	successCount, ok := device["success_count"].(int)
	if !ok {
		successCount = 0
	}
	device["risk_level"] = s.calculateRiskLevel(successCount, loginCount+1)
}

func (s *FraudDomainService) calculateRiskLevel(successCount, loginCount int) string {
	successRate := float64(successCount) / float64(loginCount)
	if successRate > 0.8 {
		return "low"
	} else if successRate > 0.5 {
		return "medium"
	}
	return "high"
}

func (s *FraudDomainService) convertDeviceMapToSlice(deviceMap map[string]map[string]interface{}) []map[string]interface{} {
	devices := make([]map[string]interface{}, 0, len(deviceMap))
	for _, device := range deviceMap {
		devices = append(devices, device)
	}
	return devices
}

func (s *FraudDomainService) TrustDevice(ctx context.Context, fingerprint string) error {
	deviceFingerprint, err := s.deviceFingerprintRepo.GetByFingerprint(ctx, fingerprint)
	if err != nil {
		return fmt.Errorf("failed to get device fingerprint: %w", err)
	}

	deviceFingerprint.Trust()
	return s.deviceFingerprintRepo.Update(ctx, deviceFingerprint)
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

func (s *FraudDomainService) GetFraudStats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"fraud_detection_enabled": true,
		"total_security_events":   0,
		"blocked_ips":             0,
		"failed_login_attempts":   0,
	}

	_, totalSecurityEvents, err := s.securityEventRepo.List(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}
	stats["total_security_events"] = totalSecurityEvents

	_, totalBlackedIPs, err := s.ipBlacklistRepo.List(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklisted IPs: %w", err)
	}
	stats["blocked_ips"] = totalBlackedIPs

	_, totalLoginAttempts, err := s.loginAttemptRepo.List(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get login attempts: %w", err)
	}
	stats["total_login_attempts"] = totalLoginAttempts

	_, highRiskEvents, err := s.securityEventRepo.GetBySeverity(ctx, "HIGH", 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get high risk events: %w", err)
	}
	stats["high_risk_events"] = highRiskEvents

	_, mediumRiskEvents, err := s.securityEventRepo.GetBySeverity(ctx, "MEDIUM", 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get medium risk events: %w", err)
	}
	stats["medium_risk_events"] = mediumRiskEvents

	activeRules, err := s.rateLimitRuleRepo.GetActiveRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active rules: %w", err)
	}
	stats["active_rate_limit_rules"] = len(activeRules)

	return stats, nil
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
