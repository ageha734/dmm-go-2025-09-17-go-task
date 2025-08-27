package usecase

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/application/dto"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
)

type FraudUsecase struct {
	fraudDomainService *service.FraudDomainService
	allowPrivateIPs    bool
	allowedPrivateIPs  map[string]bool
}

func NewFraudUsecase(fraudDomainService *service.FraudDomainService) FraudUsecaseInterface {
	return &FraudUsecase{
		fraudDomainService: fraudDomainService,
		allowPrivateIPs:    true,
		allowedPrivateIPs:  make(map[string]bool),
	}
}

func (u *FraudUsecase) AddIPToBlacklist(ctx context.Context, ip string, reason string, adminID uint) error {
	if err := u.validateIPAddress(ip); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("reason cannot be empty")
	}

	if len(reason) > 500 {
		return fmt.Errorf("reason too long (max 500 characters)")
	}

	_, err := u.fraudDomainService.GetBlacklistedIPs(ctx, 1, 1000)
	if err != nil {
		return fmt.Errorf("failed to check existing blacklist: %w", err)
	}

	clientIP := "127.0.0.1"
	userAgent := "Admin-Panel"

	return u.fraudDomainService.AddIPToBlacklist(ctx, ip, reason, clientIP, userAgent)
}

func (u *FraudUsecase) RemoveIPFromBlacklist(ctx context.Context, ip string) error {
	if err := u.validateIPAddress(ip); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	clientIP := "127.0.0.1"
	userAgent := "Admin-Panel"

	return u.fraudDomainService.RemoveIPFromBlacklist(ctx, ip, clientIP, userAgent)
}

func (u *FraudUsecase) GetBlacklistedIPs(ctx context.Context) ([]*entity.IPBlacklist, error) {
	result, err := u.fraudDomainService.GetBlacklistedIPs(ctx, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get blacklisted IPs: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from domain service")
	}

	ips, ok := resultMap["ips"].([]*entity.IPBlacklist)
	if !ok {
		return nil, fmt.Errorf("unexpected IPs type in result")
	}

	return ips, nil
}

func (u *FraudUsecase) GetSecurityEvents(ctx context.Context, limit, offset int) ([]*entity.SecurityEvent, error) {
	if limit <= 0 || limit > 1000 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	page := (offset / limit) + 1
	result, err := u.fraudDomainService.GetSecurityEvents(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from domain service")
	}

	events, ok := resultMap["events"].([]*entity.SecurityEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected events type in result")
	}

	return events, nil
}

func (u *FraudUsecase) CreateSecurityEvent(ctx context.Context, req *dto.CreateSecurityEventRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if err := u.validateSecurityEventRequest(req); err != nil {
		return fmt.Errorf("invalid security event request: %w", err)
	}

	severity := u.determineSeverity(req.EventType)

	return u.fraudDomainService.CreateSecurityEvent(
		ctx,
		req.UserID,
		req.EventType,
		req.Description,
		req.IPAddress,
		req.UserAgent,
		severity,
	)
}

func (u *FraudUsecase) CreateRateLimitRule(ctx context.Context, req *dto.CreateRateLimitRuleRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if err := u.validateRateLimitRuleRequest(req); err != nil {
		return fmt.Errorf("invalid rate limit rule request: %w", err)
	}

	return u.fraudDomainService.CreateRateLimitRule(ctx, req.RuleType, req.Identifier, req.MaxRequests, req.WindowSize)
}

func (u *FraudUsecase) UpdateRateLimitRule(ctx context.Context, id uint, req *dto.UpdateRateLimitRuleRequest) error {
	if id == 0 {
		return fmt.Errorf("invalid rule ID")
	}

	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if err := u.validateUpdateRateLimitRuleRequest(req); err != nil {
		return fmt.Errorf("invalid update rate limit rule request: %w", err)
	}

	return u.fraudDomainService.UpdateRateLimitRule(ctx, id, req.RuleType, req.Identifier, req.MaxRequests, req.WindowSize)
}

func (u *FraudUsecase) DeleteRateLimitRule(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid rule ID")
	}

	return u.fraudDomainService.DeleteRateLimitRule(ctx, id)
}

func (u *FraudUsecase) GetRateLimitRules(ctx context.Context) ([]*entity.RateLimitRule, error) {
	result, err := u.fraudDomainService.GetRateLimitRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limit rules: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from domain service")
	}

	rules, ok := resultMap["rules"].([]*entity.RateLimitRule)
	if !ok {
		return nil, fmt.Errorf("unexpected rules type in result")
	}

	return rules, nil
}

func (u *FraudUsecase) GetActiveSessions(ctx context.Context) ([]*entity.UserSession, error) {
	result, err := u.fraudDomainService.GetActiveSessions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from domain service")
	}

	sessionsData, ok := resultMap["sessions"].([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected sessions type in result")
	}

	sessions := make([]*entity.UserSession, 0, len(sessionsData))
	for range sessionsData {
		sessions = append(sessions, &entity.UserSession{})
	}

	return sessions, nil
}

func (u *FraudUsecase) DeactivateSession(ctx context.Context, sessionID string) error {
	if strings.TrimSpace(sessionID) == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	return u.fraudDomainService.DeactivateSession(ctx, sessionID)
}

func (u *FraudUsecase) GetDevices(ctx context.Context) ([]*entity.DeviceFingerprint, error) {
	result, err := u.fraudDomainService.GetDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from domain service")
	}

	devicesData, ok := resultMap["devices"].([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected devices type in result")
	}

	devices := make([]*entity.DeviceFingerprint, 0, len(devicesData))
	for range devicesData {
		devices = append(devices, &entity.DeviceFingerprint{})
	}

	return devices, nil
}

func (u *FraudUsecase) TrustDevice(ctx context.Context, fingerprint string) error {
	if strings.TrimSpace(fingerprint) == "" {
		return fmt.Errorf("fingerprint cannot be empty")
	}

	return u.fraudDomainService.TrustDevice(ctx, fingerprint)
}

func (u *FraudUsecase) CleanupExpiredData(ctx context.Context) error {
	return u.fraudDomainService.CleanupExpiredData(ctx)
}

func (u *FraudUsecase) validateIPAddress(ip string) error {
	if strings.TrimSpace(ip) == "" {
		return fmt.Errorf("IP address cannot be empty")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address format")
	}

	if parsedIP.IsPrivate() {
		if allowed, exists := u.allowedPrivateIPs[ip]; exists {
			if !allowed {
				return fmt.Errorf("private IP address %s is specifically denied", ip)
			}
		} else if !u.allowPrivateIPs {
			return fmt.Errorf("private IP addresses are not allowed")
		}
	}

	return nil
}

func (u *FraudUsecase) validateSecurityEventRequest(req *dto.CreateSecurityEventRequest) error {
	if strings.TrimSpace(req.EventType) == "" {
		return fmt.Errorf("event type cannot be empty")
	}

	if strings.TrimSpace(req.Description) == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if len(req.Description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}

	validEventTypes := []string{
		"LOGIN_SUCCESS", "LOGIN_FAILED", "PASSWORD_CHANGE", "ACCOUNT_LOCKED",
		"SUSPICIOUS_ACTIVITY", "IP_BLACKLISTED", "RATE_LIMIT_EXCEEDED",
		"SESSION_CREATED", "SESSION_EXPIRED", "DEVICE_REGISTERED",
	}

	isValidEventType := false
	for _, validType := range validEventTypes {
		if req.EventType == validType {
			isValidEventType = true
			break
		}
	}

	if !isValidEventType {
		return fmt.Errorf("invalid event type: %s", req.EventType)
	}

	if req.IPAddress != "" {
		if err := u.validateIPAddress(req.IPAddress); err != nil {
			return fmt.Errorf("invalid IP address in request: %w", err)
		}
	}

	return nil
}

func (u *FraudUsecase) validateRateLimitRuleRequest(req *dto.CreateRateLimitRuleRequest) error {
	if err := u.validateRateLimitRuleBasicFields(req); err != nil {
		return err
	}

	if err := u.validateRateLimitRuleType(req.RuleType); err != nil {
		return err
	}

	return u.validateRateLimitRuleIdentifier(req.RuleType, req.Identifier)
}

func (u *FraudUsecase) validateRateLimitRuleBasicFields(req *dto.CreateRateLimitRuleRequest) error {
	if strings.TrimSpace(req.RuleType) == "" {
		return fmt.Errorf("rule type cannot be empty")
	}

	if strings.TrimSpace(req.Identifier) == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	if req.MaxRequests <= 0 {
		return fmt.Errorf("max requests must be positive")
	}

	if req.MaxRequests > 10000 {
		return fmt.Errorf("max requests too high (max 10000)")
	}

	if req.WindowSize <= 0 {
		return fmt.Errorf("window size must be positive")
	}

	if req.WindowSize > 86400 {
		return fmt.Errorf("window size too large (max 24 hours)")
	}

	return nil
}

func (u *FraudUsecase) validateRateLimitRuleType(ruleType string) error {
	validRuleTypes := []string{"IP", "USER", "ENDPOINT", "GLOBAL"}
	for _, validType := range validRuleTypes {
		if ruleType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid rule type: %s", ruleType)
}

func (u *FraudUsecase) validateRateLimitRuleIdentifier(ruleType, identifier string) error {
	if ruleType == "IP" {
		if err := u.validateIPAddress(identifier); err != nil {
			return fmt.Errorf("invalid IP identifier: %w", err)
		}
	} else if ruleType == "USER" {
		if matched, _ := regexp.MatchString(`^\d+$`, identifier); !matched {
			return fmt.Errorf("user identifier must be numeric")
		}
	}
	return nil
}

func (u *FraudUsecase) validateUpdateRateLimitRuleRequest(req *dto.UpdateRateLimitRuleRequest) error {
	if req.RuleType != "" {
		validRuleTypes := []string{"IP", "USER", "ENDPOINT", "GLOBAL"}
		isValidRuleType := false
		for _, validType := range validRuleTypes {
			if req.RuleType == validType {
				isValidRuleType = true
				break
			}
		}

		if !isValidRuleType {
			return fmt.Errorf("invalid rule type: %s", req.RuleType)
		}
	}

	if req.MaxRequests < 0 {
		return fmt.Errorf("max requests cannot be negative")
	}

	if req.MaxRequests > 10000 {
		return fmt.Errorf("max requests too high (max 10000)")
	}

	if req.WindowSize < 0 {
		return fmt.Errorf("window size cannot be negative")
	}

	if req.WindowSize > 86400 {
		return fmt.Errorf("window size too large (max 24 hours)")
	}

	return nil
}

func (u *FraudUsecase) determineSeverity(eventType string) string {
	highSeverityEvents := []string{
		"ACCOUNT_LOCKED", "SUSPICIOUS_ACTIVITY", "IP_BLACKLISTED",
		"MULTIPLE_FAILED_LOGINS", "BRUTE_FORCE_DETECTED",
	}

	mediumSeverityEvents := []string{
		"LOGIN_FAILED", "RATE_LIMIT_EXCEEDED", "PASSWORD_CHANGE",
		"DEVICE_REGISTERED", "SESSION_EXPIRED",
	}

	for _, event := range highSeverityEvents {
		if eventType == event {
			return "HIGH"
		}
	}

	for _, event := range mediumSeverityEvents {
		if eventType == event {
			return "MEDIUM"
		}
	}

	return "LOW"
}
