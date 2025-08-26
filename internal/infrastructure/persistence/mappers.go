package persistence

import (
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
)

func UserEntityToGorm(user *entity.User) *GormUser {
	return &GormUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserGormToEntity(gormUser *GormUser) *entity.User {
	return &entity.User{
		ID:        gormUser.ID,
		Name:      gormUser.Name,
		Email:     gormUser.Email,
		Age:       gormUser.Age,
		CreatedAt: gormUser.CreatedAt,
		UpdatedAt: gormUser.UpdatedAt,
	}
}

func AuthEntityToGorm(auth *entity.Auth) *GormAuth {
	return &GormAuth{
		ID:           auth.ID,
		UserID:       auth.UserID,
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash,
		IsActive:     auth.IsActive,
		LastLoginAt:  auth.LastLoginAt,
		CreatedAt:    auth.CreatedAt,
		UpdatedAt:    auth.UpdatedAt,
	}
}

func AuthGormToEntity(gormAuth *GormAuth) *entity.Auth {
	return &entity.Auth{
		ID:           gormAuth.ID,
		UserID:       gormAuth.UserID,
		Email:        gormAuth.Email,
		PasswordHash: gormAuth.PasswordHash,
		IsActive:     gormAuth.IsActive,
		LastLoginAt:  gormAuth.LastLoginAt,
		CreatedAt:    gormAuth.CreatedAt,
		UpdatedAt:    gormAuth.UpdatedAt,
	}
}

func RoleEntityToGorm(role *entity.Role) *GormRole {
	return &GormRole{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func RoleGormToEntity(gormRole *GormRole) *entity.Role {
	return &entity.Role{
		ID:          gormRole.ID,
		Name:        gormRole.Name,
		Description: gormRole.Description,
		CreatedAt:   gormRole.CreatedAt,
		UpdatedAt:   gormRole.UpdatedAt,
	}
}

func RefreshTokenEntityToGorm(token *entity.RefreshToken) *GormRefreshToken {
	return &GormRefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		IsRevoked: token.IsRevoked,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}

func RefreshTokenGormToEntity(gormToken *GormRefreshToken) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        gormToken.ID,
		UserID:    gormToken.UserID,
		Token:     gormToken.Token,
		ExpiresAt: gormToken.ExpiresAt,
		IsRevoked: gormToken.IsRevoked,
		CreatedAt: gormToken.CreatedAt,
		UpdatedAt: gormToken.UpdatedAt,
	}
}

func MembershipTierEntityToGorm(tier *entity.MembershipTier) *GormMembershipTier {
	return &GormMembershipTier{
		ID:           tier.ID,
		Name:         tier.Name,
		Level:        tier.Level,
		Description:  tier.Description,
		Benefits:     tier.Benefits,
		Requirements: tier.Requirements,
		IsActive:     tier.IsActive,
		CreatedAt:    tier.CreatedAt,
		UpdatedAt:    tier.UpdatedAt,
	}
}

func MembershipTierGormToEntity(gormTier *GormMembershipTier) *entity.MembershipTier {
	return &entity.MembershipTier{
		ID:           gormTier.ID,
		Name:         gormTier.Name,
		Level:        gormTier.Level,
		Description:  gormTier.Description,
		Benefits:     gormTier.Benefits,
		Requirements: gormTier.Requirements,
		IsActive:     gormTier.IsActive,
		CreatedAt:    gormTier.CreatedAt,
		UpdatedAt:    gormTier.UpdatedAt,
	}
}

func UserMembershipEntityToGorm(membership *entity.UserMembership) *GormUserMembership {
	return &GormUserMembership{
		ID:             membership.ID,
		UserID:         membership.UserID,
		TierID:         membership.TierID,
		Points:         membership.Points,
		TotalSpent:     membership.TotalSpent,
		JoinedAt:       membership.JoinedAt,
		LastActivityAt: membership.LastActivityAt,
		ExpiresAt:      membership.ExpiresAt,
		IsActive:       membership.IsActive,
		CreatedAt:      membership.CreatedAt,
		UpdatedAt:      membership.UpdatedAt,
	}
}

func UserMembershipGormToEntity(gormMembership *GormUserMembership) *entity.UserMembership {
	return &entity.UserMembership{
		ID:             gormMembership.ID,
		UserID:         gormMembership.UserID,
		TierID:         gormMembership.TierID,
		Points:         gormMembership.Points,
		TotalSpent:     gormMembership.TotalSpent,
		JoinedAt:       gormMembership.JoinedAt,
		LastActivityAt: gormMembership.LastActivityAt,
		ExpiresAt:      gormMembership.ExpiresAt,
		IsActive:       gormMembership.IsActive,
		CreatedAt:      gormMembership.CreatedAt,
		UpdatedAt:      gormMembership.UpdatedAt,
	}
}

func UserProfileEntityToGorm(profile *entity.UserProfile) *GormUserProfile {
	return &GormUserProfile{
		ID:          profile.ID,
		UserID:      profile.UserID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		PhoneNumber: profile.PhoneNumber,
		DateOfBirth: profile.DateOfBirth,
		Gender:      profile.Gender,
		Address:     profile.Address,
		Preferences: profile.Preferences,
		Avatar:      profile.Avatar,
		Bio:         profile.Bio,
		IsVerified:  profile.IsVerified,
		VerifiedAt:  profile.VerifiedAt,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}

func UserProfileGormToEntity(gormProfile *GormUserProfile) *entity.UserProfile {
	return &entity.UserProfile{
		ID:          gormProfile.ID,
		UserID:      gormProfile.UserID,
		FirstName:   gormProfile.FirstName,
		LastName:    gormProfile.LastName,
		PhoneNumber: gormProfile.PhoneNumber,
		DateOfBirth: gormProfile.DateOfBirth,
		Gender:      gormProfile.Gender,
		Address:     gormProfile.Address,
		Preferences: gormProfile.Preferences,
		Avatar:      gormProfile.Avatar,
		Bio:         gormProfile.Bio,
		IsVerified:  gormProfile.IsVerified,
		VerifiedAt:  gormProfile.VerifiedAt,
		CreatedAt:   gormProfile.CreatedAt,
		UpdatedAt:   gormProfile.UpdatedAt,
	}
}

func NotificationEntityToGorm(notification *entity.Notification) *GormNotification {
	return &GormNotification{
		ID:        notification.ID,
		UserID:    notification.UserID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		Data:      notification.Data,
		IsRead:    notification.IsRead,
		ReadAt:    notification.ReadAt,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}
}

func NotificationGormToEntity(gormNotification *GormNotification) *entity.Notification {
	return &entity.Notification{
		ID:        gormNotification.ID,
		UserID:    gormNotification.UserID,
		Type:      gormNotification.Type,
		Title:     gormNotification.Title,
		Message:   gormNotification.Message,
		Data:      gormNotification.Data,
		IsRead:    gormNotification.IsRead,
		ReadAt:    gormNotification.ReadAt,
		CreatedAt: gormNotification.CreatedAt,
		UpdatedAt: gormNotification.UpdatedAt,
	}
}

func SecurityEventEntityToGorm(event *entity.SecurityEvent) *GormSecurityEvent {
	return &GormSecurityEvent{
		ID:          event.ID,
		UserID:      event.UserID,
		EventType:   event.EventType,
		Description: event.Description,
		IPAddress:   event.IPAddress,
		UserAgent:   event.UserAgent,
		Severity:    event.Severity,
		Metadata:    event.Metadata,
		CreatedAt:   event.CreatedAt,
	}
}

func SecurityEventGormToEntity(gormEvent *GormSecurityEvent) *entity.SecurityEvent {
	return &entity.SecurityEvent{
		ID:          gormEvent.ID,
		UserID:      gormEvent.UserID,
		EventType:   gormEvent.EventType,
		Description: gormEvent.Description,
		IPAddress:   gormEvent.IPAddress,
		UserAgent:   gormEvent.UserAgent,
		Severity:    gormEvent.Severity,
		Metadata:    gormEvent.Metadata,
		CreatedAt:   gormEvent.CreatedAt,
	}
}

func IPBlacklistEntityToGorm(blacklist *entity.IPBlacklist) *GormIPBlacklist {
	return &GormIPBlacklist{
		ID:        blacklist.ID,
		IPAddress: blacklist.IPAddress,
		Reason:    blacklist.Reason,
		ExpiresAt: blacklist.ExpiresAt,
		IsActive:  blacklist.IsActive,
		CreatedAt: blacklist.CreatedAt,
		UpdatedAt: blacklist.UpdatedAt,
	}
}

func IPBlacklistGormToEntity(gormBlacklist *GormIPBlacklist) *entity.IPBlacklist {
	return &entity.IPBlacklist{
		ID:        gormBlacklist.ID,
		IPAddress: gormBlacklist.IPAddress,
		Reason:    gormBlacklist.Reason,
		ExpiresAt: gormBlacklist.ExpiresAt,
		IsActive:  gormBlacklist.IsActive,
		CreatedAt: gormBlacklist.CreatedAt,
		UpdatedAt: gormBlacklist.UpdatedAt,
	}
}

func LoginAttemptEntityToGorm(attempt *entity.LoginAttempt) *GormLoginAttempt {
	return &GormLoginAttempt{
		ID:         attempt.ID,
		Email:      attempt.Email,
		IPAddress:  attempt.IPAddress,
		UserAgent:  attempt.UserAgent,
		Success:    attempt.Success,
		FailReason: attempt.FailReason,
		CreatedAt:  attempt.CreatedAt,
	}
}

func LoginAttemptGormToEntity(gormAttempt *GormLoginAttempt) *entity.LoginAttempt {
	return &entity.LoginAttempt{
		ID:         gormAttempt.ID,
		Email:      gormAttempt.Email,
		IPAddress:  gormAttempt.IPAddress,
		UserAgent:  gormAttempt.UserAgent,
		Success:    gormAttempt.Success,
		FailReason: gormAttempt.FailReason,
		CreatedAt:  gormAttempt.CreatedAt,
	}
}

func RateLimitRuleEntityToGorm(rule *entity.RateLimitRule) *GormRateLimitRule {
	return &GormRateLimitRule{
		ID:          rule.ID,
		Name:        rule.Name,
		Resource:    rule.Resource,
		MaxRequests: rule.MaxRequests,
		WindowSize:  rule.WindowSize,
		IsActive:    rule.IsActive,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}
}

func RateLimitRuleGormToEntity(gormRule *GormRateLimitRule) *entity.RateLimitRule {
	return &entity.RateLimitRule{
		ID:          gormRule.ID,
		Name:        gormRule.Name,
		Resource:    gormRule.Resource,
		MaxRequests: gormRule.MaxRequests,
		WindowSize:  gormRule.WindowSize,
		IsActive:    gormRule.IsActive,
		CreatedAt:   gormRule.CreatedAt,
		UpdatedAt:   gormRule.UpdatedAt,
	}
}

func UserSessionEntityToGorm(session *entity.UserSession) *GormUserSession {
	return &GormUserSession{
		ID:        session.ID,
		UserID:    session.UserID,
		SessionID: session.SessionID,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		ExpiresAt: session.ExpiresAt,
		IsActive:  session.IsActive,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}
}

func UserSessionGormToEntity(gormSession *GormUserSession) *entity.UserSession {
	return &entity.UserSession{
		ID:        gormSession.ID,
		UserID:    gormSession.UserID,
		SessionID: gormSession.SessionID,
		IPAddress: gormSession.IPAddress,
		UserAgent: gormSession.UserAgent,
		ExpiresAt: gormSession.ExpiresAt,
		IsActive:  gormSession.IsActive,
		CreatedAt: gormSession.CreatedAt,
		UpdatedAt: gormSession.UpdatedAt,
	}
}

func DeviceFingerprintEntityToGorm(fingerprint *entity.DeviceFingerprint) *GormDeviceFingerprint {
	return &GormDeviceFingerprint{
		ID:          fingerprint.ID,
		UserID:      fingerprint.UserID,
		Fingerprint: fingerprint.Fingerprint,
		DeviceInfo:  fingerprint.DeviceInfo,
		IsTrusted:   fingerprint.IsTrusted,
		LastSeenAt:  fingerprint.LastSeenAt,
		CreatedAt:   fingerprint.CreatedAt,
		UpdatedAt:   fingerprint.UpdatedAt,
	}
}

func DeviceFingerprintGormToEntity(gormFingerprint *GormDeviceFingerprint) *entity.DeviceFingerprint {
	return &entity.DeviceFingerprint{
		ID:          gormFingerprint.ID,
		UserID:      gormFingerprint.UserID,
		Fingerprint: gormFingerprint.Fingerprint,
		DeviceInfo:  gormFingerprint.DeviceInfo,
		IsTrusted:   gormFingerprint.IsTrusted,
		LastSeenAt:  gormFingerprint.LastSeenAt,
		CreatedAt:   gormFingerprint.CreatedAt,
		UpdatedAt:   gormFingerprint.UpdatedAt,
	}
}
