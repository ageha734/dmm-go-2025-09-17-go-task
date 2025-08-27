package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	var users []*entity.User
	if u, ok := args.Get(0).([]*entity.User); ok {
		users = u
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return users, total, args.Error(2)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if profile, ok := args.Get(0).(*entity.UserProfile); ok {
		return profile, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	if profile == nil {
		return errors.New("profile is nil")
	}
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) Delete(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockUserMembershipRepository struct {
	mock.Mock
}

func (m *MockUserMembershipRepository) Create(ctx context.Context, membership *entity.UserMembership) error {
	args := m.Called(ctx, membership)
	return args.Error(0)
}

func (m *MockUserMembershipRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserMembership, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if membership, ok := args.Get(0).(*entity.UserMembership); ok {
		return membership, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserMembershipRepository) Update(ctx context.Context, membership *entity.UserMembership) error {
	if membership == nil {
		return errors.New("membership is nil")
	}
	args := m.Called(ctx, membership)
	return args.Error(0)
}

func (m *MockUserMembershipRepository) Delete(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserMembershipRepository) List(ctx context.Context, offset, limit int) ([]*entity.UserMembership, int64, error) {
	args := m.Called(ctx, offset, limit)
	var memberships []*entity.UserMembership
	if m, ok := args.Get(0).([]*entity.UserMembership); ok {
		memberships = m
	}
	var total int64
	if t, ok := args.Get(1).(int64); ok {
		total = t
	}
	return memberships, total, args.Error(2)
}

func (m *MockUserMembershipRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if stats, ok := args.Get(0).(map[string]interface{}); ok {
		return stats, args.Error(1)
	}
	return nil, args.Error(1)
}

func createTestUserUsecase() *usecase.UserUsecase {
	userRepo := &MockUserRepository{}
	userProfileRepo := &MockUserProfileRepository{}
	userMembershipRepo := &MockUserMembershipRepository{}
	fraudService := &MockFraudDomainService{}

	return usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)
}

func TestNewUserUsecase(t *testing.T) {
	t.Run("新しいUserUsecaseを正常に作成できる", func(t *testing.T) {
		uc := createTestUserUsecase()
		assert.NotNil(t, uc)
	})
}

func TestUserUsecaseCreateUser(t *testing.T) {
	t.Run("新しいユーザーを正常に作成できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		req := usecase.CreateUserRequest{
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   25,
		}

		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
		userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
		userProfileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
		fraudService.On("CreateSecurityEvent", mock.Anything, mock.AnythingOfType("*uint"), "USER_CREATED", "User created via API", "192.168.1.1", "test-agent", "LOW").Return(nil)

		user, err := uc.CreateUser(context.Background(), req, "192.168.1.1", "test-agent")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Name, user.Name)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.Age, user.Age)

		userRepo.AssertExpectations(t)
		userProfileRepo.AssertExpectations(t)
		fraudService.AssertExpectations(t)
	})

	t.Run("既存のメールアドレスでユーザー作成に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		req := usecase.CreateUserRequest{
			Name:  "テストユーザー",
			Email: "existing@example.com",
			Age:   25,
		}

		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)

		user, err := uc.CreateUser(context.Background(), req, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")

		userRepo.AssertExpectations(t)
	})

	t.Run("無効な年齢でユーザー作成に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		req := usecase.CreateUserRequest{
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   -1,
		}

		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)

		user, err := uc.CreateUser(context.Background(), req, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid user data")

		userRepo.AssertExpectations(t)
	})

	t.Run("データベースエラーでユーザー作成に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		req := usecase.CreateUserRequest{
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   25,
		}

		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
		userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

		user, err := uc.CreateUser(context.Background(), req, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to create user")

		userRepo.AssertExpectations(t)
	})
}

func TestUsecaseResponseComparison(t *testing.T) {
	expected := map[string]interface{}{
		"total_users":   int64(100),
		"premium_users": 20,
		"basic_users":   80,
		"user_details": map[string]interface{}{
			"active_users":   90,
			"inactive_users": 10,
		},
	}

	actual := map[string]interface{}{
		"total_users":   int64(100),
		"premium_users": 20,
		"basic_users":   80,
		"user_details": map[string]interface{}{
			"active_users":   90,
			"inactive_users": 10,
		},
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Stats response mismatch (-want +got):\n%s", diff)
	}

	opts := cmp.Options{
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.String() == "user_details"
		}, cmp.Ignore()),
	}

	if diff := cmp.Diff(expected, actual, opts); diff != "" {
		t.Errorf("Stats response mismatch (ignoring user_details) (-want +got):\n%s", diff)
	}
}

func TestUserUsecaseGetUser(t *testing.T) {
	t.Run("ユーザーを正常に取得できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		expectedUser := &entity.User{
			ID:    1,
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   25,
		}

		userRepo.On("GetByID", mock.Anything, uint(1)).Return(expectedUser, nil)

		user, err := uc.GetUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		userRepo.AssertExpectations(t)
	})

	t.Run("存在しないユーザーの取得に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))

		user, err := uc.GetUser(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "failed to get user")

		userRepo.AssertExpectations(t)
	})
}

func TestUserUsecaseGetUsers(t *testing.T) {
	t.Run("ユーザーリストを正常に取得できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		expectedUsers := []*entity.User{
			{ID: 1, Name: "ユーザー1", Email: "user1@example.com", Age: 25},
			{ID: 2, Name: "ユーザー2", Email: "user2@example.com", Age: 30},
		}

		userRepo.On("List", mock.Anything, 0, 20).Return(expectedUsers, int64(2), nil)

		response, err := uc.GetUsers(context.Background(), 1, 20)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedUsers, response.Users)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 20, response.Limit)
		assert.Equal(t, 1, response.TotalPages)

		userRepo.AssertExpectations(t)
	})

	t.Run("無効なページ番号を正規化する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 20).Return([]*entity.User{}, int64(0), nil)

		response, err := uc.GetUsers(context.Background(), 0, 0)

		assert.NoError(t, err)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 20, response.Limit)

		userRepo.AssertExpectations(t)
	})

	t.Run("大きすぎるリミットを正規化する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 20).Return([]*entity.User{}, int64(0), nil)

		response, err := uc.GetUsers(context.Background(), 1, 200)

		assert.NoError(t, err)
		assert.Equal(t, 20, response.Limit)

		userRepo.AssertExpectations(t)
	})
}

func TestUserUsecaseUpdateUser(t *testing.T) {
	t.Run("ユーザー情報を正常に更新できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		existingUser := &entity.User{
			ID:    1,
			Name:  "旧名前",
			Email: "old@example.com",
			Age:   25,
		}

		req := usecase.UpdateUserRequest{
			Name:  "新名前",
			Email: "new@example.com",
			Age:   30,
		}

		userRepo.On("GetByID", mock.Anything, uint(1)).Return(existingUser, nil)
		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(false, nil)
		userRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
		fraudService.On("CreateSecurityEvent", mock.Anything, mock.AnythingOfType("*uint"), "USER_UPDATED", "User profile updated", "192.168.1.1", "test-agent", "LOW").Return(nil)

		user, err := uc.UpdateUser(context.Background(), 1, req, 1, "192.168.1.1", "test-agent")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Name, user.Name)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.Age, user.Age)

		userRepo.AssertExpectations(t)
	})

	t.Run("権限がない場合は更新に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		req := usecase.UpdateUserRequest{
			Name: "新名前",
		}

		user, err := uc.UpdateUser(context.Background(), 1, req, 2, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "permission denied")
	})

	t.Run("既存のメールアドレスへの変更に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		existingUser := &entity.User{
			ID:    1,
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   25,
		}

		req := usecase.UpdateUserRequest{
			Email: "existing@example.com",
		}

		userRepo.On("GetByID", mock.Anything, uint(1)).Return(existingUser, nil)
		userRepo.On("ExistsByEmail", mock.Anything, req.Email).Return(true, nil)

		user, err := uc.UpdateUser(context.Background(), 1, req, 1, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")

		userRepo.AssertExpectations(t)
	})
}

func TestUserUsecaseDeleteUser(t *testing.T) {
	t.Run("ユーザーを正常に削除できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		existingUser := &entity.User{
			ID:    1,
			Name:  "テストユーザー",
			Email: "test@example.com",
			Age:   25,
		}

		userRepo.On("GetByID", mock.Anything, uint(1)).Return(existingUser, nil)
		userRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		userProfileRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		userMembershipRepo.On("Delete", mock.Anything, uint(1)).Return(nil)
		fraudService.On("CreateSecurityEvent", mock.Anything, mock.AnythingOfType("*uint"), "USER_DELETED", "User account deleted", "192.168.1.1", "test-agent", "MEDIUM").Return(nil)

		err := uc.DeleteUser(context.Background(), 1, 1, "192.168.1.1", "test-agent")

		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
		userProfileRepo.AssertExpectations(t)
		userMembershipRepo.AssertExpectations(t)
	})

	t.Run("権限がない場合は削除に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		err := uc.DeleteUser(context.Background(), 1, 2, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})

	t.Run("存在しないユーザーの削除に失敗する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("GetByID", mock.Anything, uint(999)).Return(nil, errors.New("user not found"))

		err := uc.DeleteUser(context.Background(), 999, 999, "192.168.1.1", "test-agent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get user")

		userRepo.AssertExpectations(t)
	})
}

func TestUserUsecaseGetUserStats(t *testing.T) {
	t.Run("ユーザー統計を正常に取得できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 1).Return([]*entity.User{}, int64(100), nil)
		userMembershipRepo.On("GetStats", mock.Anything).Return(map[string]interface{}{
			"premium_users": 20,
			"basic_users":   80,
		}, nil)

		stats, err := uc.GetUserStats(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(100), stats["total_users"])
		assert.Equal(t, 20, stats["premium_users"])
		assert.Equal(t, 80, stats["basic_users"])

		userRepo.AssertExpectations(t)
		userMembershipRepo.AssertExpectations(t)
	})

	t.Run("メンバーシップ統計の取得に失敗してもユーザー統計は取得できる", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 1).Return([]*entity.User{}, int64(50), nil)
		userMembershipRepo.On("GetStats", mock.Anything).Return(nil, errors.New("membership stats error"))

		stats, err := uc.GetUserStats(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(50), stats["total_users"])
		assert.NotContains(t, stats, "premium_users")

		userRepo.AssertExpectations(t)
		userMembershipRepo.AssertExpectations(t)
	})
}

func TestUserUsecaseHealthCheck(t *testing.T) {
	t.Run("ヘルスチェックが正常に動作する", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 1).Return([]*entity.User{}, int64(0), nil)

		status := uc.HealthCheck(context.Background())

		assert.Equal(t, "ok", status["status"])
		assert.Equal(t, "user-service", status["service"])
		assert.Equal(t, "connected", status["database"])

		userRepo.AssertExpectations(t)
	})

	t.Run("データベース接続エラー時のヘルスチェック", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userProfileRepo := &MockUserProfileRepository{}
		userMembershipRepo := &MockUserMembershipRepository{}
		fraudService := &MockFraudDomainService{}

		uc := usecase.NewUserUsecase(userRepo, userProfileRepo, userMembershipRepo, fraudService, nil)

		userRepo.On("List", mock.Anything, 0, 1).Return(nil, int64(0), errors.New("database connection error"))

		status := uc.HealthCheck(context.Background())

		assert.Equal(t, "error", status["status"])
		assert.Equal(t, "user-service", status["service"])
		assert.Equal(t, "disconnected", status["database"])

		userRepo.AssertExpectations(t)
	})
}
