package persistence

import (
	"context"
	"fmt"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repository.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Create(ctx context.Context, auth *entity.Auth) error {
	gormAuth := AuthEntityToGorm(auth)
	if err := r.db.WithContext(ctx).Create(gormAuth).Error; err != nil {
		return err
	}
	auth.ID = gormAuth.ID
	return nil
}

func (r *authRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Auth, error) {
	var gormAuth GormAuth
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&gormAuth).Error; err != nil {
		return nil, err
	}
	return AuthGormToEntity(&gormAuth), nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*entity.Auth, error) {
	var gormAuth GormAuth
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gormAuth).Error; err != nil {
		return nil, err
	}
	return AuthGormToEntity(&gormAuth), nil
}

func (r *authRepository) Update(ctx context.Context, auth *entity.Auth) error {
	gormAuth := AuthEntityToGorm(auth)
	return r.db.WithContext(ctx).Save(gormAuth).Error
}

func (r *authRepository) Delete(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&GormAuth{}).Error
}

func (r *authRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&GormAuth{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	gormRole := RoleEntityToGorm(role)
	if err := r.db.WithContext(ctx).Create(gormRole).Error; err != nil {
		return err
	}
	role.ID = gormRole.ID
	return nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	var gormRole GormRole
	if err := r.db.WithContext(ctx).First(&gormRole, id).Error; err != nil {
		return nil, err
	}
	return RoleGormToEntity(&gormRole), nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	var gormRole GormRole
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&gormRole).Error; err != nil {
		return nil, err
	}
	return RoleGormToEntity(&gormRole), nil
}

func (r *roleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	var gormRoles []GormRole
	if err := r.db.WithContext(ctx).Find(&gormRoles).Error; err != nil {
		return nil, err
	}

	roles := make([]*entity.Role, len(gormRoles))
	for i, gormRole := range gormRoles {
		roles[i] = RoleGormToEntity(&gormRole)
	}

	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	gormRole := RoleEntityToGorm(role)
	return r.db.WithContext(ctx).Save(gormRole).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&GormRole{}, id).Error
}

func (r *roleRepository) AssignToUser(ctx context.Context, userID, roleID uint) error {
	userRole := &GormUserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *roleRepository) RemoveFromUser(ctx context.Context, userID, roleID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&GormUserRole{}).Error
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error) {
	var gormRoles []GormRole
	if err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&gormRoles).Error; err != nil {
		return nil, err
	}

	roles := make([]*entity.Role, len(gormRoles))
	for i, gormRole := range gormRoles {
		roles[i] = RoleGormToEntity(&gormRole)
	}

	return roles, nil
}

func (r *roleRepository) GetUserRoleNames(ctx context.Context, userID uint) ([]string, error) {
	var roleNames []string
	if err := r.db.WithContext(ctx).
		Model(&GormRole{}).
		Select("roles.name").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("name", &roleNames).Error; err != nil {
		return nil, err
	}

	if len(roleNames) == 0 {
		return nil, fmt.Errorf("user %d has no roles assigned and default roles do not exist", userID)
	}

	return roleNames, nil
}

func (r *roleRepository) AssignAdminRole(ctx context.Context, userID uint) error {
	var adminRole GormRole
	if err := r.db.WithContext(ctx).Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&GormUserRole{}).
		Where("user_id = ? AND role_id = ?", userID, adminRole.ID).
		Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		userRoleAssignment := &GormUserRole{
			UserID: userID,
			RoleID: adminRole.ID,
		}
		return r.db.WithContext(ctx).Create(userRoleAssignment).Error
	}

	return nil
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) repository.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	gormToken := RefreshTokenEntityToGorm(token)
	if err := r.db.WithContext(ctx).Create(gormToken).Error; err != nil {
		return err
	}
	token.ID = gormToken.ID
	return nil
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	var gormToken GormRefreshToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&gormToken).Error; err != nil {
		return nil, err
	}
	return RefreshTokenGormToEntity(&gormToken), nil
}

func (r *refreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	gormToken := RefreshTokenEntityToGorm(token)
	return r.db.WithContext(ctx).Save(gormToken).Error
}

func (r *refreshTokenRepository) RevokeByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&GormRefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < NOW()").Delete(&GormRefreshToken{}).Error
}
