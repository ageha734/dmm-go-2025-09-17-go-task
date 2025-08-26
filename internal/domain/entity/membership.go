package entity

import (
	"errors"
	"time"
)

var (
	ErrInsufficientPoints = errors.New("insufficient points")
	ErrInvalidPointAmount = errors.New("invalid point amount")
)

type MembershipTier struct {
	ID           uint
	Name         string
	Level        int
	Description  string
	Benefits     string
	Requirements string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewMembershipTier(name string, level int, description, benefits, requirements string) *MembershipTier {
	now := time.Now()
	return &MembershipTier{
		Name:         name,
		Level:        level,
		Description:  description,
		Benefits:     benefits,
		Requirements: requirements,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type UserMembership struct {
	ID             uint
	UserID         uint
	TierID         uint
	Points         int
	TotalSpent     float64
	JoinedAt       time.Time
	LastActivityAt *time.Time
	ExpiresAt      *time.Time
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUserMembership(userID, tierID uint) *UserMembership {
	now := time.Now()
	return &UserMembership{
		UserID:     userID,
		TierID:     tierID,
		Points:     0,
		TotalSpent: 0,
		JoinedAt:   now,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (um *UserMembership) AddPoints(points int) error {
	if points <= 0 {
		return ErrInvalidPointAmount
	}
	um.Points += points
	um.UpdateLastActivity()
	return nil
}

func (um *UserMembership) SpendPoints(points int) error {
	if points <= 0 {
		return ErrInvalidPointAmount
	}
	if um.Points < points {
		return ErrInsufficientPoints
	}
	um.Points -= points
	um.UpdateLastActivity()
	return nil
}

func (um *UserMembership) UpdateTotalSpent(amount float64) {
	if amount > 0 {
		um.TotalSpent += amount
		um.UpdateLastActivity()
	}
}

func (um *UserMembership) UpdateLastActivity() {
	now := time.Now()
	um.LastActivityAt = &now
	um.UpdatedAt = now
}

func (um *UserMembership) Deactivate() {
	um.IsActive = false
	um.UpdatedAt = time.Now()
}

type PointTransaction struct {
	ID            uint
	UserID        uint
	Type          string
	Points        int
	Description   string
	ReferenceType string
	ReferenceID   *uint
	ExpiresAt     *time.Time
	CreatedAt     time.Time
}

func NewPointTransaction(userID uint, transactionType string, points int, description string) *PointTransaction {
	return &PointTransaction{
		UserID:      userID,
		Type:        transactionType,
		Points:      points,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

type UserProfile struct {
	ID          uint
	UserID      uint
	FirstName   string
	LastName    string
	PhoneNumber string
	DateOfBirth *time.Time
	Gender      string
	Address     *string
	Preferences *string
	Avatar      string
	Bio         string
	IsVerified  bool
	VerifiedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUserProfile(userID uint) *UserProfile {
	now := time.Now()
	return &UserProfile{
		UserID:     userID,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (up *UserProfile) UpdateProfile(firstName, lastName, phoneNumber, gender, bio string, dateOfBirth *time.Time) {
	if firstName != "" {
		up.FirstName = firstName
	}
	if lastName != "" {
		up.LastName = lastName
	}
	if phoneNumber != "" {
		up.PhoneNumber = phoneNumber
	}
	if gender != "" {
		up.Gender = gender
	}
	if bio != "" {
		up.Bio = bio
	}
	if dateOfBirth != nil {
		up.DateOfBirth = dateOfBirth
	}
	up.UpdatedAt = time.Now()
}

func (up *UserProfile) Verify() {
	now := time.Now()
	up.IsVerified = true
	up.VerifiedAt = &now
	up.UpdatedAt = now
}

type Notification struct {
	ID        uint
	UserID    uint
	Type      string
	Title     string
	Message   string
	Data      *string
	IsRead    bool
	ReadAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNotification(userID uint, notificationType, title, message string, data *string) *Notification {
	now := time.Now()
	return &Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     title,
		Message:   message,
		Data:      data,
		IsRead:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n *Notification) MarkAsRead() {
	if !n.IsRead {
		now := time.Now()
		n.IsRead = true
		n.ReadAt = &now
		n.UpdatedAt = now
	}
}
