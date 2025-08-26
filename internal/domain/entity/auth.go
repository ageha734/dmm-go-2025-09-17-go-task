package entity

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrWeakPassword    = errors.New("password is too weak")
)

type Auth struct {
	ID           uint
	UserID       uint
	Email        string
	PasswordHash string
	IsActive     bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAuth(userID uint, email, password string) (*Auth, error) {
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Auth{
		UserID:       userID,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (a *Auth) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
}

func (a *Auth) ChangePassword(currentPassword, newPassword string) error {
	if err := a.VerifyPassword(currentPassword); err != nil {
		return ErrInvalidPassword
	}

	if err := validatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.PasswordHash = string(hashedPassword)
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Auth) UpdateLastLogin() {
	now := time.Now()
	a.LastLoginAt = &now
	a.UpdatedAt = now
}

func (a *Auth) Deactivate() {
	a.IsActive = false
	a.UpdatedAt = time.Now()
}

func (a *Auth) Activate() {
	a.IsActive = true
	a.UpdatedAt = time.Now()
}

// ValidatePassword validates if a password meets the minimum requirements
func (a *Auth) ValidatePassword(password string) error {
	return validatePassword(password)
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return ErrWeakPassword
	}
	return nil
}

type Role struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewRole(name, description string) *Role {
	now := time.Now()
	return &Role{
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type RefreshToken struct {
	ID        uint
	UserID    uint
	Token     string
	ExpiresAt time.Time
	IsRevoked bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewRefreshToken(userID uint, token string, expiresAt time.Time) *RefreshToken {
	now := time.Now()
	return &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
	rt.UpdatedAt = time.Now()
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsRevoked && !rt.IsExpired()
}
