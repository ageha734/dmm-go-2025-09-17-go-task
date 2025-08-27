package persistence

import (
	"context"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"gorm.io/gorm"
)

type securityEventRepository struct {
	db *gorm.DB
}

func NewSecurityEventRepository(db *gorm.DB) repository.SecurityEventRepository {
	return &securityEventRepository{db: db}
}

func (r *securityEventRepository) Create(ctx context.Context, event *entity.SecurityEvent) error {
	gormEvent := SecurityEventEntityToGorm(event)
	if err := r.db.WithContext(ctx).Create(gormEvent).Error; err != nil {
		return err
	}
	event.ID = gormEvent.ID
	return nil
}

func (r *securityEventRepository) GetByID(ctx context.Context, id uint) (*entity.SecurityEvent, error) {
	var gormEvent GormSecurityEvent
	if err := r.db.WithContext(ctx).First(&gormEvent, id).Error; err != nil {
		return nil, err
	}
	return SecurityEventGormToEntity(&gormEvent), nil
}

func (r *securityEventRepository) List(ctx context.Context, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	var gormEvents []GormSecurityEvent
	var total int64

	if err := r.db.WithContext(ctx).Model(&GormSecurityEvent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&gormEvents).Error; err != nil {
		return nil, 0, err
	}

	events := make([]*entity.SecurityEvent, len(gormEvents))
	for i, gormEvent := range gormEvents {
		events[i] = SecurityEventGormToEntity(&gormEvent)
	}

	return events, total, nil
}

func (r *securityEventRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	var gormEvents []GormSecurityEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&GormSecurityEvent{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&gormEvents).Error; err != nil {
		return nil, 0, err
	}

	events := make([]*entity.SecurityEvent, len(gormEvents))
	for i, gormEvent := range gormEvents {
		events[i] = SecurityEventGormToEntity(&gormEvent)
	}

	return events, total, nil
}

func (r *securityEventRepository) GetBySeverity(ctx context.Context, severity string, offset, limit int) ([]*entity.SecurityEvent, int64, error) {
	var gormEvents []GormSecurityEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&GormSecurityEvent{}).Where("severity = ?", severity)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&gormEvents).Error; err != nil {
		return nil, 0, err
	}

	events := make([]*entity.SecurityEvent, len(gormEvents))
	for i, gormEvent := range gormEvents {
		events[i] = SecurityEventGormToEntity(&gormEvent)
	}

	return events, total, nil
}

func (r *securityEventRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&GormSecurityEvent{}, id).Error
}

type ipBlacklistRepository struct {
	db *gorm.DB
}

func NewIPBlacklistRepository(db *gorm.DB) repository.IPBlacklistRepository {
	return &ipBlacklistRepository{db: db}
}

func (r *ipBlacklistRepository) Create(ctx context.Context, blacklist *entity.IPBlacklist) error {
	gormBlacklist := IPBlacklistEntityToGorm(blacklist)
	if err := r.db.WithContext(ctx).Create(gormBlacklist).Error; err != nil {
		return err
	}
	blacklist.ID = gormBlacklist.ID
	return nil
}

func (r *ipBlacklistRepository) GetByIP(ctx context.Context, ipAddress string) (*entity.IPBlacklist, error) {
	var gormBlacklist GormIPBlacklist
	if err := r.db.WithContext(ctx).Where("ip_address = ?", ipAddress).First(&gormBlacklist).Error; err != nil {
		return nil, err
	}
	return IPBlacklistGormToEntity(&gormBlacklist), nil
}

func (r *ipBlacklistRepository) List(ctx context.Context, offset, limit int) ([]*entity.IPBlacklist, int64, error) {
	var gormBlacklists []GormIPBlacklist
	var total int64

	if err := r.db.WithContext(ctx).Model(&GormIPBlacklist{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&gormBlacklists).Error; err != nil {
		return nil, 0, err
	}

	blacklists := make([]*entity.IPBlacklist, len(gormBlacklists))
	for i, gormBlacklist := range gormBlacklists {
		blacklists[i] = IPBlacklistGormToEntity(&gormBlacklist)
	}

	return blacklists, total, nil
}

func (r *ipBlacklistRepository) Update(ctx context.Context, blacklist *entity.IPBlacklist) error {
	gormBlacklist := IPBlacklistEntityToGorm(blacklist)
	return r.db.WithContext(ctx).Save(gormBlacklist).Error
}

func (r *ipBlacklistRepository) Delete(ctx context.Context, ipAddress string) error {
	return r.db.WithContext(ctx).Where("ip_address = ?", ipAddress).Delete(&GormIPBlacklist{}).Error
}

func (r *ipBlacklistRepository) IsBlacklisted(ctx context.Context, ipAddress string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&GormIPBlacklist{}).
		Where("ip_address = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > NOW())", ipAddress, true).
		Count(&count).Error
	return count > 0, err
}

func (r *ipBlacklistRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&GormIPBlacklist{}).
		Where("expires_at IS NOT NULL AND expires_at <= NOW()").
		Update("is_active", false).Error
}

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(db *gorm.DB) repository.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) Create(ctx context.Context, attempt *entity.LoginAttempt) error {
	gormAttempt := LoginAttemptEntityToGorm(attempt)
	if err := r.db.WithContext(ctx).Create(gormAttempt).Error; err != nil {
		return err
	}
	attempt.ID = gormAttempt.ID
	return nil
}

func (r *loginAttemptRepository) GetByEmail(ctx context.Context, email string, since time.Time) ([]*entity.LoginAttempt, error) {
	var gormAttempts []GormLoginAttempt
	if err := r.db.WithContext(ctx).
		Where("email = ? AND created_at >= ?", email, since).
		Order("created_at DESC").
		Find(&gormAttempts).Error; err != nil {
		return nil, err
	}

	attempts := make([]*entity.LoginAttempt, len(gormAttempts))
	for i, gormAttempt := range gormAttempts {
		attempts[i] = LoginAttemptGormToEntity(&gormAttempt)
	}

	return attempts, nil
}

func (r *loginAttemptRepository) GetByIP(ctx context.Context, ipAddress string, since time.Time) ([]*entity.LoginAttempt, error) {
	var gormAttempts []GormLoginAttempt
	if err := r.db.WithContext(ctx).
		Where("ip_address = ? AND created_at >= ?", ipAddress, since).
		Order("created_at DESC").
		Find(&gormAttempts).Error; err != nil {
		return nil, err
	}

	attempts := make([]*entity.LoginAttempt, len(gormAttempts))
	for i, gormAttempt := range gormAttempts {
		attempts[i] = LoginAttemptGormToEntity(&gormAttempt)
	}

	return attempts, nil
}

func (r *loginAttemptRepository) CountFailedAttempts(ctx context.Context, email string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&GormLoginAttempt{}).
		Where("email = ? AND success = ? AND created_at >= ?", email, false, since).
		Count(&count).Error
	return count, err
}

func (r *loginAttemptRepository) List(ctx context.Context, offset, limit int) ([]*entity.LoginAttempt, int64, error) {
	var gormAttempts []GormLoginAttempt
	var total int64

	if err := r.db.WithContext(ctx).Model(&GormLoginAttempt{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&gormAttempts).Error; err != nil {
		return nil, 0, err
	}

	attempts := make([]*entity.LoginAttempt, len(gormAttempts))
	for i, gormAttempt := range gormAttempts {
		attempts[i] = LoginAttemptGormToEntity(&gormAttempt)
	}

	return attempts, total, nil
}

func (r *loginAttemptRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&GormLoginAttempt{}, id).Error
}

func (r *loginAttemptRepository) CleanupOld(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Where("created_at < ?", before).Delete(&GormLoginAttempt{}).Error
}
