package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
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
	user, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user, ok := args.Get(0).(*entity.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
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

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	var users []*entity.User
	if v := args.Get(0); v != nil {
		u, ok := v.([]*entity.User)
		if !ok {
			return nil, 0, args.Error(2)
		}
		users = u
	}
	var total int64
	if tv, ok := args.Get(1).(int64); ok {
		total = tv
	}
	return users, total, args.Error(2)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) Create(ctx context.Context, auth *entity.Auth) error {
	args := m.Called(ctx, auth)
	return args.Error(0)
}

func (m *MockAuthRepository) GetByUserID(ctx context.Context, userID uint) (*entity.Auth, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if auth, ok := args.Get(0).(*entity.Auth); ok {
		return auth, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepository) GetByEmail(ctx context.Context, email string) (*entity.Auth, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if auth, ok := args.Get(0).(*entity.Auth); ok {
		return auth, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthRepository) Update(ctx context.Context, auth *entity.Auth) error {
	if auth == nil {
		return errors.New("auth is nil")
	}
	args := m.Called(ctx, auth)
	return args.Error(0)
}

func (m *MockAuthRepository) Delete(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if role, ok := args.Get(0).(*entity.Role); ok {
		return role, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if role, ok := args.Get(0).(*entity.Role); ok {
		return role, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]*entity.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if roles, ok := args.Get(0).([]*entity.Role); ok {
		return roles, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) AssignToUser(ctx context.Context, userID, roleID uint) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRoleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if roles, ok := args.Get(0).([]*entity.Role); ok {
		return roles, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	if role == nil {
		return errors.New("role is nil")
	}
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveFromUser(ctx context.Context, userID, roleID uint) error {
	if userID == 0 || roleID == 0 {
		return errors.New("invalid userID or roleID")
	}
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if refreshToken, ok := args.Get(0).(*entity.RefreshToken); ok {
		return refreshToken, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRefreshTokenRepository) Update(ctx context.Context, token *entity.RefreshToken) error {
	if token == nil {
		return errors.New("refresh token is nil")
	}
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByUserID(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestAuthDomainServiceRegister(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		userName  string
		age       int
		setupMock func(*MockUserRepository, *MockAuthRepository, *MockRoleRepository)
		wantErr   bool
	}{
		{
			name:     "正常なユーザー登録",
			email:    "test@example.com",
			password: "password123",
			userName: "テストユーザー",
			age:      25,
			setupMock: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				ctx := context.Background()
				authRepo.On("ExistsByEmail", ctx, "test@example.com").Return(false, nil)
				userRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
				authRepo.On("Create", ctx, mock.AnythingOfType("*entity.Auth")).Return(nil)
				roleRepo.On("GetByName", ctx, "user").Return(entity.NewRole("user", "Default user role"), nil)
				roleRepo.On("AssignToUser", ctx, mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "既存のメールアドレスでエラー",
			email:    "existing@example.com",
			password: "password123",
			userName: "テストユーザー",
			age:      25,
			setupMock: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				ctx := context.Background()
				authRepo.On("ExistsByEmail", ctx, "existing@example.com").Return(true, nil)
			},
			wantErr: true,
		},
		{
			name:     "弱いパスワードでエラー",
			email:    "test@example.com",
			password: "123",
			userName: "テストユーザー",
			age:      25,
			setupMock: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				tempAuth := &entity.Auth{}
				err := tempAuth.ValidatePassword("123")
				assert.Error(t, err, "Password validation should fail for weak password")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			roleRepo := new(MockRoleRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			tt.setupMock(userRepo, authRepo, roleRepo)

			service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

			ctx := context.Background()
			user, err := service.Register(ctx, tt.userName, tt.email, tt.password, tt.age)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.userName, user.Name)
				assert.Equal(t, tt.age, user.Age)
			}

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
		})
	}
}

func TestAuthDomainServiceLogin(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		setupMock func(*MockAuthRepository, *MockRoleRepository)
		wantErr   bool
	}{
		{
			name:     "正常なログイン",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				ctx := context.Background()
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				roles := []*entity.Role{entity.NewRole("user", "Default user role")}

				authRepo.On("GetByEmail", ctx, "test@example.com").Return(auth, nil)
				authRepo.On("Update", ctx, mock.AnythingOfType("*entity.Auth")).Return(nil)
				roleRepo.On("GetUserRoles", ctx, uint(1)).Return(roles, nil)
			},
			wantErr: false,
		},
		{
			name:     "存在しないユーザー",
			email:    "nonexistent@example.com",
			password: "password123",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				ctx := context.Background()
				authRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return((*entity.Auth)(nil), assert.AnError)
			},
			wantErr: true,
		},
		{
			name:     "間違ったパスワード",
			email:    "test@example.com",
			password: "wrong-password",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository) {
				ctx := context.Background()
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				authRepo.On("GetByEmail", ctx, "test@example.com").Return(auth, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			roleRepo := new(MockRoleRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			tt.setupMock(authRepo, roleRepo)

			service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

			ctx := context.Background()
			auth, roles, err := service.Login(ctx, tt.email, tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, auth)
				assert.Nil(t, roles)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auth)
				assert.NotNil(t, roles)
			}

			authRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
		})
	}
}

func TestAuthDomainServiceGenerateAccessToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	authRepo := new(MockAuthRepository)
	roleRepo := new(MockRoleRepository)
	refreshTokenRepo := new(MockRefreshTokenRepository)
	service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

	userID := uint(1)
	email := "test@example.com"
	roles := []string{"user"}

	token, err := service.GenerateAccessToken(userID, email, roles)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthDomainServiceValidateToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	authRepo := new(MockAuthRepository)
	roleRepo := new(MockRoleRepository)
	refreshTokenRepo := new(MockRefreshTokenRepository)
	service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

	userID := uint(1)
	email := "test@example.com"
	roles := []string{"user"}

	validToken, err := service.GenerateAccessToken(userID, email, roles)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "有効なトークン",
			token:   validToken,
			wantErr: false,
		},
		{
			name:    "無効なトークン",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "空のトークン",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.UserID)
				assert.Equal(t, email, claims.Email)
			}
		})
	}
}

func TestAuthDomainServiceGenerateRefreshToken(t *testing.T) {
	userRepo := new(MockUserRepository)
	authRepo := new(MockAuthRepository)
	roleRepo := new(MockRoleRepository)
	refreshTokenRepo := new(MockRefreshTokenRepository)

	ctx := context.Background()
	refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)

	service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

	userID := uint(1)
	token, err := service.GenerateRefreshToken(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	refreshTokenRepo.AssertExpectations(t)
}

func TestAuthDomainServiceRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		setupMock func(*MockAuthRepository, *MockRoleRepository, *MockRefreshTokenRepository)
		wantErr   bool
	}{
		{
			name:  "有効なリフレッシュトークン",
			token: "valid-refresh-token",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository, refreshTokenRepo *MockRefreshTokenRepository) {
				ctx := context.Background()
				refreshToken := entity.NewRefreshToken(1, "valid-refresh-token", time.Now().Add(24*time.Hour))
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				roles := []*entity.Role{entity.NewRole("user", "Default user role")}

				refreshTokenRepo.On("GetByToken", ctx, "valid-refresh-token").Return(refreshToken, nil)
				authRepo.On("GetByUserID", ctx, uint(1)).Return(auth, nil)
				roleRepo.On("GetUserRoles", ctx, uint(1)).Return(roles, nil)
				refreshTokenRepo.On("Update", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				refreshTokenRepo.On("Create", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "存在しないリフレッシュトークン",
			token: "nonexistent-token",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository, refreshTokenRepo *MockRefreshTokenRepository) {
				ctx := context.Background()
				refreshTokenRepo.On("GetByToken", ctx, "nonexistent-token").Return((*entity.RefreshToken)(nil), assert.AnError)
			},
			wantErr: true,
		},
		{
			name:  "期限切れのリフレッシュトークン",
			token: "expired-token",
			setupMock: func(authRepo *MockAuthRepository, roleRepo *MockRoleRepository, refreshTokenRepo *MockRefreshTokenRepository) {
				ctx := context.Background()
				refreshToken := entity.NewRefreshToken(1, "expired-token", time.Now().Add(-1*time.Hour))
				refreshTokenRepo.On("GetByToken", ctx, "expired-token").Return(refreshToken, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			roleRepo := new(MockRoleRepository)
			refreshTokenRepo := new(MockRefreshTokenRepository)
			tt.setupMock(authRepo, roleRepo, refreshTokenRepo)

			service := service.NewAuthDomainService(userRepo, authRepo, roleRepo, refreshTokenRepo, "test-secret")

			ctx := context.Background()
			auth, roles, newToken, err := service.RefreshToken(ctx, tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, auth)
				assert.Nil(t, roles)
				assert.Empty(t, newToken)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auth)
				assert.NotNil(t, roles)
				assert.NotEmpty(t, newToken)
			}

			authRepo.AssertExpectations(t)
			roleRepo.AssertExpectations(t)
			refreshTokenRepo.AssertExpectations(t)
		})
	}
}
