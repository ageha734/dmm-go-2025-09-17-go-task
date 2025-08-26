package persistence

import (
	"context"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	gormUser := UserEntityToGorm(user)
	if err := r.db.WithContext(ctx).Create(gormUser).Error; err != nil {
		return err
	}
	user.ID = gormUser.ID
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var gormUser GormUser
	if err := r.db.WithContext(ctx).First(&gormUser, id).Error; err != nil {
		return nil, err
	}
	return UserGormToEntity(&gormUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var gormUser GormUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser).Error; err != nil {
		return nil, err
	}
	return UserGormToEntity(&gormUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	gormUser := UserEntityToGorm(user)
	return r.db.WithContext(ctx).Save(gormUser).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&GormUser{}, id).Error
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	var gormUsers []GormUser
	var total int64

	if err := r.db.WithContext(ctx).Model(&GormUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&gormUsers).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*entity.User, len(gormUsers))
	for i, gormUser := range gormUsers {
		users[i] = UserGormToEntity(&gormUser)
	}

	return users, total, nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&GormUser{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
