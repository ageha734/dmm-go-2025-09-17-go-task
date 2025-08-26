package persistence

import (
	"time"

	"gorm.io/gorm"
)

type GormUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Age       int            `json:"age"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormUser) TableName() string {
	return "users"
}

type GormAuth struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormAuth) TableName() string {
	return "auths"
}

type GormRole struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormRole) TableName() string {
	return "roles"
}

type GormUserRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	RoleID    uint           `json:"role_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
	Role GormRole `json:"role" gorm:"foreignKey:RoleID"`
}

func (GormUserRole) TableName() string {
	return "user_roles"
}

type GormRefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at"`
	IsRevoked bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormRefreshToken) TableName() string {
	return "refresh_tokens"
}

type GormMembershipTier struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"uniqueIndex;not null"`
	Level        int            `json:"level" gorm:"not null;index"`
	Description  string         `json:"description"`
	Benefits     string         `json:"benefits" gorm:"type:json"`
	Requirements string         `json:"requirements" gorm:"type:json"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormMembershipTier) TableName() string {
	return "membership_tiers"
}

type GormUserMembership struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	TierID         uint           `json:"tier_id" gorm:"not null"`
	Points         int            `json:"points" gorm:"default:0"`
	TotalSpent     float64        `json:"total_spent" gorm:"default:0"`
	JoinedAt       time.Time      `json:"joined_at"`
	LastActivityAt *time.Time     `json:"last_activity_at"`
	ExpiresAt      *time.Time     `json:"expires_at"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser           `json:"user" gorm:"foreignKey:UserID"`
	Tier GormMembershipTier `json:"tier" gorm:"foreignKey:TierID"`
}

func (GormUserMembership) TableName() string {
	return "user_memberships"
}

type GormUserProfile struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	FirstName   string         `json:"first_name" gorm:"default:''"`
	LastName    string         `json:"last_name" gorm:"default:''"`
	PhoneNumber string         `json:"phone_number" gorm:"default:''"`
	DateOfBirth *time.Time     `json:"date_of_birth"`
	Gender      string         `json:"gender" gorm:"default:''"`
	Address     *string        `json:"address" gorm:"type:json;default:NULL"`
	Preferences *string        `json:"preferences" gorm:"type:json;default:NULL"`
	Avatar      string         `json:"avatar" gorm:"default:''"`
	Bio         string         `json:"bio" gorm:"default:''"`
	IsVerified  bool           `json:"is_verified" gorm:"default:false"`
	VerifiedAt  *time.Time     `json:"verified_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormUserProfile) TableName() string {
	return "user_profiles"
}

type GormNotification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Type      string         `json:"type" gorm:"not null;index"`
	Title     string         `json:"title" gorm:"not null"`
	Message   string         `json:"message" gorm:"not null"`
	Data      *string        `json:"data" gorm:"type:json"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormNotification) TableName() string {
	return "notifications"
}

type GormSecurityEvent struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      *uint          `json:"user_id" gorm:"index"`
	EventType   string         `json:"event_type" gorm:"not null;index"`
	Description string         `json:"description"`
	IPAddress   string         `json:"ip_address" gorm:"not null;index"`
	UserAgent   string         `json:"user_agent"`
	Severity    string         `json:"severity" gorm:"not null;index"`
	Metadata    *string        `json:"metadata" gorm:"type:json"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User *GormUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (GormSecurityEvent) TableName() string {
	return "security_events"
}

type GormIPBlacklist struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	IPAddress string         `json:"ip_address" gorm:"uniqueIndex;not null"`
	Reason    string         `json:"reason"`
	ExpiresAt *time.Time     `json:"expires_at"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormIPBlacklist) TableName() string {
	return "ip_blacklists"
}

type GormLoginAttempt struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Email      string         `json:"email" gorm:"not null;index"`
	IPAddress  string         `json:"ip_address" gorm:"not null;index"`
	UserAgent  string         `json:"user_agent"`
	Success    bool           `json:"success"`
	FailReason string         `json:"fail_reason"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormLoginAttempt) TableName() string {
	return "login_attempts"
}

type GormRateLimitRule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Resource    string         `json:"resource" gorm:"not null"`
	MaxRequests int            `json:"max_requests" gorm:"not null"`
	WindowSize  int            `json:"window_size" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (GormRateLimitRule) TableName() string {
	return "rate_limit_rules"
}

type GormUserSession struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	SessionID string         `json:"session_id" gorm:"uniqueIndex;not null"`
	IPAddress string         `json:"ip_address" gorm:"not null"`
	UserAgent string         `json:"user_agent"`
	ExpiresAt time.Time      `json:"expires_at"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormUserSession) TableName() string {
	return "user_sessions"
}

type GormDeviceFingerprint struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Fingerprint string         `json:"fingerprint" gorm:"not null;index"`
	DeviceInfo  *string        `json:"device_info" gorm:"type:json"`
	IsTrusted   bool           `json:"is_trusted" gorm:"default:false"`
	LastSeenAt  time.Time      `json:"last_seen_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User GormUser `json:"user" gorm:"foreignKey:UserID"`
}

func (GormDeviceFingerprint) TableName() string {
	return "device_fingerprints"
}
