package entity

import (
	"errors"
	"time"
)

var (
	ErrIPBlacklisted      = errors.New("IP address is blacklisted")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrSuspiciousActivity = errors.New("suspicious activity detected")
)

type SecurityEvent struct {
	ID          uint
	UserID      *uint
	EventType   string
	Description string
	IPAddress   string
	UserAgent   string
	Severity    string
	Metadata    *string
	CreatedAt   time.Time
}

func NewSecurityEvent(userID *uint, eventType, description, ipAddress, userAgent, severity string) *SecurityEvent {
	return &SecurityEvent{
		UserID:      userID,
		EventType:   eventType,
		Description: description,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Severity:    severity,
		CreatedAt:   time.Now(),
	}
}

func (se *SecurityEvent) IsHighSeverity() bool {
	return se.Severity == "HIGH" || se.Severity == "CRITICAL"
}

type IPBlacklist struct {
	ID        uint
	IPAddress string
	Reason    string
	ExpiresAt *time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewIPBlacklist(ipAddress, reason string, expiresAt *time.Time) *IPBlacklist {
	now := time.Now()
	return &IPBlacklist{
		IPAddress: ipAddress,
		Reason:    reason,
		ExpiresAt: expiresAt,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (ib *IPBlacklist) IsExpired() bool {
	if ib.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ib.ExpiresAt)
}

func (ib *IPBlacklist) Deactivate() {
	ib.IsActive = false
	ib.UpdatedAt = time.Now()
}

type LoginAttempt struct {
	ID         uint
	Email      string
	IPAddress  string
	UserAgent  string
	Success    bool
	FailReason string
	CreatedAt  time.Time
}

func NewLoginAttempt(email, ipAddress, userAgent string, success bool, failReason string) *LoginAttempt {
	return &LoginAttempt{
		Email:      email,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Success:    success,
		FailReason: failReason,
		CreatedAt:  time.Now(),
	}
}

type RateLimitRule struct {
	ID          uint
	Name        string
	Resource    string
	MaxRequests int
	WindowSize  int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewRateLimitRule(name, resource string, maxRequests, windowSize int) *RateLimitRule {
	now := time.Now()
	return &RateLimitRule{
		Name:        name,
		Resource:    resource,
		MaxRequests: maxRequests,
		WindowSize:  windowSize,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (rlr *RateLimitRule) UpdateRule(maxRequests, windowSize int) {
	if maxRequests > 0 {
		rlr.MaxRequests = maxRequests
	}
	if windowSize > 0 {
		rlr.WindowSize = windowSize
	}
	rlr.UpdatedAt = time.Now()
}

func (rlr *RateLimitRule) Deactivate() {
	rlr.IsActive = false
	rlr.UpdatedAt = time.Now()
}

type UserSession struct {
	ID        uint
	UserID    uint
	SessionID string
	IPAddress string
	UserAgent string
	ExpiresAt time.Time
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUserSession(userID uint, sessionID, ipAddress, userAgent string, expiresAt time.Time) *UserSession {
	now := time.Now()
	return &UserSession{
		UserID:    userID,
		SessionID: sessionID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: expiresAt,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (us *UserSession) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}

func (us *UserSession) Deactivate() {
	us.IsActive = false
	us.UpdatedAt = time.Now()
}

func (us *UserSession) IsValid() bool {
	return us.IsActive && !us.IsExpired()
}

type DeviceFingerprint struct {
	ID          uint
	UserID      uint
	Fingerprint string
	DeviceInfo  *string
	IsTrusted   bool
	LastSeenAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewDeviceFingerprint(userID uint, fingerprint string, deviceInfo *string) *DeviceFingerprint {
	now := time.Now()
	return &DeviceFingerprint{
		UserID:      userID,
		Fingerprint: fingerprint,
		DeviceInfo:  deviceInfo,
		IsTrusted:   false,
		LastSeenAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (df *DeviceFingerprint) Trust() {
	df.IsTrusted = true
	df.UpdatedAt = time.Now()
}

func (df *DeviceFingerprint) UpdateLastSeen() {
	df.LastSeenAt = time.Now()
	df.UpdatedAt = time.Now()
}

type FraudAnalysis struct {
	RiskScore      float64
	RiskLevel      string
	Factors        []string
	Recommendation string
}

func NewFraudAnalysis(riskScore float64, factors []string) *FraudAnalysis {
	var riskLevel, recommendation string

	switch {
	case riskScore >= 0.8:
		riskLevel = "HIGH"
		recommendation = "Block request and require additional verification"
	case riskScore >= 0.5:
		riskLevel = "MEDIUM"
		recommendation = "Monitor closely and consider additional verification"
	default:
		riskLevel = "LOW"
		recommendation = "Allow request"
	}

	return &FraudAnalysis{
		RiskScore:      riskScore,
		RiskLevel:      riskLevel,
		Factors:        factors,
		Recommendation: recommendation,
	}
}

func (fa *FraudAnalysis) IsHighRisk() bool {
	return fa.RiskLevel == "HIGH"
}
