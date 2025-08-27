package usecase

import (
	"context"
	"fmt"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/repository"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
)

type UserUsecase struct {
	userRepo           repository.UserRepository
	userProfileRepo    repository.UserProfileRepository
	userMembershipRepo repository.UserMembershipRepository
	fraudDomainService service.FraudDomainServiceInterface
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"min=0,max=150"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"omitempty,email"`
	Age   int    `json:"age" binding:"min=0,max=150"`
}

type UserListResponse struct {
	Users      []*entity.User `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	userProfileRepo repository.UserProfileRepository,
	userMembershipRepo repository.UserMembershipRepository,
	fraudDomainService service.FraudDomainServiceInterface,
) *UserUsecase {
	return &UserUsecase{
		userRepo:           userRepo,
		userProfileRepo:    userProfileRepo,
		userMembershipRepo: userMembershipRepo,
		fraudDomainService: fraudDomainService,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req CreateUserRequest, ipAddress, userAgent string) (*entity.User, error) {
	exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := entity.NewUser(req.Name, req.Email, req.Age)
	if !user.IsValidAge() || !user.IsValidEmail() {
		return nil, fmt.Errorf("invalid user data")
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	profile := entity.NewUserProfile(user.ID)
	_ = u.userProfileRepo.Create(ctx, profile)

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &user.ID, "USER_CREATED",
		"User created via API", ipAddress, userAgent, "LOW")

	return user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (u *UserUsecase) GetUserProfile(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return user, nil
}

func (u *UserUsecase) GetUsers(ctx context.Context, page, limit int) (*UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, total, err := u.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &UserListResponse{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, userID uint, req UpdateUserRequest, requestUserID uint, ipAddress, userAgent string) (*entity.User, error) {
	if userID != requestUserID {
		return nil, fmt.Errorf("permission denied")
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.Email != "" && req.Email != user.Email {
		exists, err := u.userRepo.ExistsByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
		user.Email = req.Email
	}

	user.UpdateProfile(req.Name, req.Age)

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &userID, "USER_UPDATED",
		"User profile updated", ipAddress, userAgent, "LOW")

	return user, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, userID uint, requestUserID uint, ipAddress, userAgent string) error {
	if userID != requestUserID {
		return fmt.Errorf("permission denied")
	}

	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := u.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	_ = u.userProfileRepo.Delete(ctx, userID)

	_ = u.userMembershipRepo.Delete(ctx, userID)

	_ = u.fraudDomainService.CreateSecurityEvent(ctx, &userID, "USER_DELETED",
		"User account deleted", ipAddress, userAgent, "MEDIUM")

	return nil
}

func (u *UserUsecase) GetUserStats(ctx context.Context) (map[string]interface{}, error) {
	_, total, err := u.userRepo.List(ctx, 0, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	stats := map[string]interface{}{
		"total_users": total,
	}

	membershipStats, err := u.userMembershipRepo.GetStats(ctx)
	if err == nil {
		for k, v := range membershipStats {
			stats[k] = v
		}
	}

	return stats, nil
}

func (u *UserUsecase) GetFraudStats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"fraud_detection_enabled": true,
		"total_security_events":   0,
		"blocked_ips":             0,
		"failed_login_attempts":   0,
	}

	fraudStats, err := u.fraudDomainService.GetFraudStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud stats: %w", err)
	}

	for k, v := range fraudStats {
		stats[k] = v
	}

	return stats, nil
}

func (u *UserUsecase) HealthCheck(ctx context.Context) map[string]string {
	status := map[string]string{
		"status":  "ok",
		"service": "user-service",
	}

	_, _, err := u.userRepo.List(ctx, 0, 1)
	if err != nil {
		status["status"] = "error"
		status["database"] = "disconnected"
	} else {
		status["database"] = "connected"
	}

	return status
}
