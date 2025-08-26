package entity_test

import (
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewSecurityEvent(t *testing.T) {
	tests := []struct {
		name        string
		userID      *uint
		eventType   string
		description string
		ipAddress   string
		userAgent   string
		severity    string
	}{
		{
			name:        "正常なセキュリティイベント作成",
			userID:      &[]uint{1}[0],
			eventType:   "LOGIN",
			description: "User logged in successfully",
			ipAddress:   "192.168.1.1",
			userAgent:   "Mozilla/5.0",
			severity:    "LOW",
		},
		{
			name:        "ユーザーIDなしのセキュリティイベント",
			userID:      nil,
			eventType:   "FAILED_LOGIN",
			description: "Failed login attempt",
			ipAddress:   "192.168.1.100",
			userAgent:   "curl/7.68.0",
			severity:    "HIGH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := entity.NewSecurityEvent(tt.userID, tt.eventType, tt.description, tt.ipAddress, tt.userAgent, tt.severity)

			assert.NotNil(t, event)
			assert.Equal(t, tt.userID, event.UserID)
			assert.Equal(t, tt.eventType, event.EventType)
			assert.Equal(t, tt.description, event.Description)
			assert.Equal(t, tt.ipAddress, event.IPAddress)
			assert.Equal(t, tt.userAgent, event.UserAgent)
			assert.Equal(t, tt.severity, event.Severity)
			assert.False(t, event.CreatedAt.IsZero())
		})
	}
}

func TestSecurityEventIsHighSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity string
		want     bool
	}{
		{
			name:     "HIGH重要度",
			severity: "HIGH",
			want:     true,
		},
		{
			name:     "CRITICAL重要度",
			severity: "CRITICAL",
			want:     true,
		},
		{
			name:     "MEDIUM重要度",
			severity: "MEDIUM",
			want:     false,
		},
		{
			name:     "LOW重要度",
			severity: "LOW",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := entity.NewSecurityEvent(nil, "TEST", "Test event", "127.0.0.1", "test", tt.severity)
			got := event.IsHighSeverity()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewIPBlacklist(t *testing.T) {
	expiresAt := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name      string
		ipAddress string
		reason    string
		expiresAt *time.Time
	}{
		{
			name:      "期限付きIPブラックリスト",
			ipAddress: "192.168.1.100",
			reason:    "Multiple failed login attempts",
			expiresAt: &expiresAt,
		},
		{
			name:      "永続IPブラックリスト",
			ipAddress: "10.0.0.1",
			reason:    "Malicious activity",
			expiresAt: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blacklist := entity.NewIPBlacklist(tt.ipAddress, tt.reason, tt.expiresAt)

			assert.NotNil(t, blacklist)
			assert.Equal(t, tt.ipAddress, blacklist.IPAddress)
			assert.Equal(t, tt.reason, blacklist.Reason)
			assert.Equal(t, tt.expiresAt, blacklist.ExpiresAt)
			assert.True(t, blacklist.IsActive)
			assert.False(t, blacklist.CreatedAt.IsZero())
			assert.False(t, blacklist.UpdatedAt.IsZero())
		})
	}
}

func TestIPBlacklistIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt *time.Time
		want      bool
	}{
		{
			name:      "未来の期限で有効",
			expiresAt: &[]time.Time{time.Now().Add(1 * time.Hour)}[0],
			want:      false,
		},
		{
			name:      "過去の期限で期限切れ",
			expiresAt: &[]time.Time{time.Now().Add(-1 * time.Hour)}[0],
			want:      true,
		},
		{
			name:      "期限なしで永続",
			expiresAt: nil,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blacklist := entity.NewIPBlacklist("192.168.1.1", "test", tt.expiresAt)
			got := blacklist.IsExpired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIPBlacklistDeactivate(t *testing.T) {
	blacklist := entity.NewIPBlacklist("192.168.1.1", "test", nil)
	assert.True(t, blacklist.IsActive)

	oldUpdatedAt := blacklist.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	blacklist.Deactivate()

	assert.False(t, blacklist.IsActive)
	assert.True(t, blacklist.UpdatedAt.After(oldUpdatedAt))
}

func TestNewLoginAttempt(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		ipAddress  string
		userAgent  string
		success    bool
		failReason string
	}{
		{
			name:       "成功したログイン試行",
			email:      "test@example.com",
			ipAddress:  "192.168.1.1",
			userAgent:  "Mozilla/5.0",
			success:    true,
			failReason: "",
		},
		{
			name:       "失敗したログイン試行",
			email:      "test@example.com",
			ipAddress:  "192.168.1.100",
			userAgent:  "curl/7.68.0",
			success:    false,
			failReason: "Invalid password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempt := entity.NewLoginAttempt(tt.email, tt.ipAddress, tt.userAgent, tt.success, tt.failReason)

			assert.NotNil(t, attempt)
			assert.Equal(t, tt.email, attempt.Email)
			assert.Equal(t, tt.ipAddress, attempt.IPAddress)
			assert.Equal(t, tt.userAgent, attempt.UserAgent)
			assert.Equal(t, tt.success, attempt.Success)
			assert.Equal(t, tt.failReason, attempt.FailReason)
			assert.False(t, attempt.CreatedAt.IsZero())
		})
	}
}

func TestNewRateLimitRule(t *testing.T) {
	rule := entity.NewRateLimitRule("login_attempts", "/api/auth/login", 5, 300)

	assert.NotNil(t, rule)
	assert.Equal(t, "login_attempts", rule.Name)
	assert.Equal(t, "/api/auth/login", rule.Resource)
	assert.Equal(t, 5, rule.MaxRequests)
	assert.Equal(t, 300, rule.WindowSize)
	assert.True(t, rule.IsActive)
	assert.False(t, rule.CreatedAt.IsZero())
	assert.False(t, rule.UpdatedAt.IsZero())
}

func TestRateLimitRuleUpdateRule(t *testing.T) {
	rule := entity.NewRateLimitRule("test", "/api/test", 10, 60)
	oldUpdatedAt := rule.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	rule.UpdateRule(20, 120)

	assert.Equal(t, 20, rule.MaxRequests)
	assert.Equal(t, 120, rule.WindowSize)
	assert.True(t, rule.UpdatedAt.After(oldUpdatedAt))
}

func TestRateLimitRuleDeactivate(t *testing.T) {
	rule := entity.NewRateLimitRule("test", "/api/test", 10, 60)
	assert.True(t, rule.IsActive)

	oldUpdatedAt := rule.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	rule.Deactivate()

	assert.False(t, rule.IsActive)
	assert.True(t, rule.UpdatedAt.After(oldUpdatedAt))
}

func TestNewUserSession(t *testing.T) {
	userID := uint(1)
	sessionID := "session123"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	expiresAt := time.Now().Add(24 * time.Hour)

	session := entity.NewUserSession(userID, sessionID, ipAddress, userAgent, expiresAt)

	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, sessionID, session.SessionID)
	assert.Equal(t, ipAddress, session.IPAddress)
	assert.Equal(t, userAgent, session.UserAgent)
	assert.Equal(t, expiresAt, session.ExpiresAt)
	assert.True(t, session.IsActive)
	assert.False(t, session.CreatedAt.IsZero())
	assert.False(t, session.UpdatedAt.IsZero())
}

func TestUserSessionIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "未来の期限で有効",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "過去の期限で期限切れ",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := entity.NewUserSession(1, "test", "127.0.0.1", "test", tt.expiresAt)
			got := session.IsExpired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserSessionDeactivate(t *testing.T) {
	session := entity.NewUserSession(1, "test", "127.0.0.1", "test", time.Now().Add(1*time.Hour))
	assert.True(t, session.IsActive)

	oldUpdatedAt := session.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	session.Deactivate()

	assert.False(t, session.IsActive)
	assert.True(t, session.UpdatedAt.After(oldUpdatedAt))
}

func TestUserSessionIsValid(t *testing.T) {
	tests := []struct {
		name      string
		isActive  bool
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "アクティブで有効期限内",
			isActive:  true,
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      true,
		},
		{
			name:      "非アクティブ",
			isActive:  false,
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "期限切れ",
			isActive:  true,
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := entity.NewUserSession(1, "test", "127.0.0.1", "test", tt.expiresAt)
			if !tt.isActive {
				session.Deactivate()
			}
			got := session.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewDeviceFingerprint(t *testing.T) {
	userID := uint(1)
	fingerprint := "fp123456"
	deviceInfo := "Chrome/91.0 on Windows 10"

	device := entity.NewDeviceFingerprint(userID, fingerprint, &deviceInfo)

	assert.NotNil(t, device)
	assert.Equal(t, userID, device.UserID)
	assert.Equal(t, fingerprint, device.Fingerprint)
	assert.Equal(t, &deviceInfo, device.DeviceInfo)
	assert.False(t, device.IsTrusted)
	assert.False(t, device.LastSeenAt.IsZero())
	assert.False(t, device.CreatedAt.IsZero())
	assert.False(t, device.UpdatedAt.IsZero())
}

func TestDeviceFingerprintTrust(t *testing.T) {
	device := entity.NewDeviceFingerprint(1, "fp123", nil)
	assert.False(t, device.IsTrusted)

	oldUpdatedAt := device.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	device.Trust()

	assert.True(t, device.IsTrusted)
	assert.True(t, device.UpdatedAt.After(oldUpdatedAt))
}

func TestDeviceFingerprintUpdateLastSeen(t *testing.T) {
	device := entity.NewDeviceFingerprint(1, "fp123", nil)
	oldLastSeenAt := device.LastSeenAt
	oldUpdatedAt := device.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	device.UpdateLastSeen()

	assert.True(t, device.LastSeenAt.After(oldLastSeenAt))
	assert.True(t, device.UpdatedAt.After(oldUpdatedAt))
}

func TestNewFraudAnalysis(t *testing.T) {
	tests := []struct {
		name          string
		riskScore     float64
		factors       []string
		expectedLevel string
		expectedRec   string
	}{
		{
			name:          "高リスク分析",
			riskScore:     0.9,
			factors:       []string{"suspicious_ip", "multiple_failed_attempts"},
			expectedLevel: "HIGH",
			expectedRec:   "Block request and require additional verification",
		},
		{
			name:          "中リスク分析",
			riskScore:     0.6,
			factors:       []string{"new_device"},
			expectedLevel: "MEDIUM",
			expectedRec:   "Monitor closely and consider additional verification",
		},
		{
			name:          "低リスク分析",
			riskScore:     0.2,
			factors:       []string{},
			expectedLevel: "LOW",
			expectedRec:   "Allow request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := entity.NewFraudAnalysis(tt.riskScore, tt.factors)

			assert.NotNil(t, analysis)
			assert.Equal(t, tt.riskScore, analysis.RiskScore)
			assert.Equal(t, tt.expectedLevel, analysis.RiskLevel)
			assert.Equal(t, tt.factors, analysis.Factors)
			assert.Equal(t, tt.expectedRec, analysis.Recommendation)
		})
	}
}

func TestFraudAnalysisIsHighRisk(t *testing.T) {
	tests := []struct {
		name      string
		riskLevel string
		want      bool
	}{
		{
			name:      "高リスク",
			riskLevel: "HIGH",
			want:      true,
		},
		{
			name:      "中リスク",
			riskLevel: "MEDIUM",
			want:      false,
		},
		{
			name:      "低リスク",
			riskLevel: "LOW",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &entity.FraudAnalysis{RiskLevel: tt.riskLevel}
			got := analysis.IsHighRisk()
			assert.Equal(t, tt.want, got)
		})
	}
}
