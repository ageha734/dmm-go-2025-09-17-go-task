package persistence

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	gormProfile := UserProfileEntityToGorm(profile)
	if err := r.db.WithContext(ctx).Create(gormProfile).Error; err != nil {
		return err
	}
	profile.ID = gormProfile.ID
	return nil
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error) {
	var gormProfile GormUserProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&gormProfile).Error; err != nil {
		return nil, err
	}
	return UserProfileGormToEntity(&gormProfile), nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	gormProfile := UserProfileEntityToGorm(profile)
	return r.db.WithContext(ctx).Save(gormProfile).Error
}

func (r *userProfileRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&GormUserProfile{}).Error
}

type userMembershipRepository struct {
	db *gorm.DB
}

func NewUserMembershipRepository(db *gorm.DB) repository.UserMembershipRepository {
	return &userMembershipRepository{db: db}
}

func (r *userMembershipRepository) Create(ctx context.Context, membership *entity.UserMembership) error {
	gormMembership := UserMembershipEntityToGorm(membership)
	if err := r.db.WithContext(ctx).Create(gormMembership).Error; err != nil {
		return err
	}
	membership.ID = gormMembership.ID
	return nil
}

func (r *userMembershipRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserMembership, error) {
	var gormMembership GormUserMembership
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&gormMembership).Error; err != nil {
		return nil, err
	}
	return UserMembershipGormToEntity(&gormMembership), nil
}

func (r *userMembershipRepository) Update(ctx context.Context, membership *entity.UserMembership) error {
	gormMembership := UserMembershipEntityToGorm(membership)
	return r.db.WithContext(ctx).Save(gormMembership).Error
}

func (r *userMembershipRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&GormUserMembership{}).Error
}

func (r *userMembershipRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	var totalMembers int64
	r.db.WithContext(ctx).Model(&GormUserMembership{}).Count(&totalMembers)

	return map[string]interface{}{
		"total_members": totalMembers,
	}, nil
}

func (r *userMembershipRepository) List(ctx context.Context, offset, limit int) ([]*entity.UserMembership, int64, error) {
	var gormMemberships []GormUserMembership
	var total int64

	if err := r.db.WithContext(ctx).Model(&GormUserMembership{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&gormMemberships).Error; err != nil {
		return nil, 0, err
	}

	memberships := make([]*entity.UserMembership, len(gormMemberships))
	for i, gormMembership := range gormMemberships {
		memberships[i] = UserMembershipGormToEntity(&gormMembership)
	}

	return memberships, total, nil
}

type rateLimitRuleRepository struct {
	db *gorm.DB
}

func NewRateLimitRuleRepository(db *gorm.DB) repository.RateLimitRuleRepository {
	return &rateLimitRuleRepository{db: db}
}

func (r *rateLimitRuleRepository) Create(ctx context.Context, rule *entity.RateLimitRule) error {
	return nil
}

func (r *rateLimitRuleRepository) GetByID(ctx context.Context, id uint) (*entity.RateLimitRule, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *rateLimitRuleRepository) GetByResource(ctx context.Context, resource string) (*entity.RateLimitRule, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *rateLimitRuleRepository) List(ctx context.Context) ([]*entity.RateLimitRule, error) {
	return nil, nil
}

func (r *rateLimitRuleRepository) Update(ctx context.Context, rule *entity.RateLimitRule) error {
	return nil
}
func (r *rateLimitRuleRepository) Delete(ctx context.Context, id uint) error { return nil }
func (r *rateLimitRuleRepository) GetActiveRules(ctx context.Context) ([]*entity.RateLimitRule, error) {
	return nil, nil
}

type userSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) repository.UserSessionRepository {
	return &userSessionRepository{db: db}
}

func (r *userSessionRepository) Create(ctx context.Context, session *entity.UserSession) error {
	return nil
}

func (r *userSessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*entity.UserSession, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *userSessionRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.UserSession, error) {
	return nil, nil
}

func (r *userSessionRepository) Update(ctx context.Context, session *entity.UserSession) error {
	return nil
}
func (r *userSessionRepository) Delete(ctx context.Context, sessionID string) error { return nil }
func (r *userSessionRepository) DeactivateByUserID(ctx context.Context, userID uint) error {
	return nil
}
func (r *userSessionRepository) CleanupExpired(ctx context.Context) error { return nil }

type deviceFingerprintRepository struct {
	db *gorm.DB
}

func NewDeviceFingerprintRepository(db *gorm.DB) repository.DeviceFingerprintRepository {
	return &deviceFingerprintRepository{db: db}
}

func (r *deviceFingerprintRepository) Create(ctx context.Context, fingerprint *entity.DeviceFingerprint) error {
	return nil
}

func (r *deviceFingerprintRepository) GetByFingerprint(ctx context.Context, fingerprint string) (*entity.DeviceFingerprint, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *deviceFingerprintRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.DeviceFingerprint, error) {
	return nil, nil
}

func (r *deviceFingerprintRepository) Update(ctx context.Context, fingerprint *entity.DeviceFingerprint) error {
	return nil
}
func (r *deviceFingerprintRepository) Delete(ctx context.Context, id uint) error { return nil }
func (r *deviceFingerprintRepository) IsTrustedDevice(ctx context.Context, userID uint, fingerprint string) (bool, error) {
	return false, nil
}
