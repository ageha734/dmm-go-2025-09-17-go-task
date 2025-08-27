package usecase_test

import (
	"context"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/service"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthDomainService struct {
	mock.Mock
}

func (m *MockAuthDomainService) Register(ctx context.Context, name, email, password string, age int) (*entity.User, error) {
	args := m.Called(ctx, name, email, password, age)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if user, ok := args.Get(0).(*entity.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthDomainService) Login(ctx context.Context, email, password string) (*entity.Auth, []string, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var auth *entity.Auth
	if a, ok := args.Get(0).(*entity.Auth); ok {
		auth = a
	}
	var roles []string
	if r, ok := args.Get(1).([]string); ok {
		roles = r
	}
	return auth, roles, args.Error(2)
}

func (m *MockAuthDomainService) GenerateAccessToken(userID uint, email string, roles []string) (string, error) {
	args := m.Called(userID, email, roles)
	return args.String(0), args.Error(1)
}

func (m *MockAuthDomainService) GenerateRefreshToken(ctx context.Context, userID uint) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthDomainService) ValidateToken(tokenString string) (*service.JWTClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if claims, ok := args.Get(0).(*service.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthDomainService) RefreshToken(ctx context.Context, refreshTokenStr string) (*entity.Auth, []string, string, error) {
	args := m.Called(ctx, refreshTokenStr)
	if args.Get(0) == nil {
		return nil, nil, "", args.Error(3)
	}
	var auth *entity.Auth
	if a, ok := args.Get(0).(*entity.Auth); ok {
		auth = a
	}
	var roles []string
	if r, ok := args.Get(1).([]string); ok {
		roles = r
	}
	return auth, roles, args.String(2), args.Error(3)
}

func (m *MockAuthDomainService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func (m *MockAuthDomainService) Logout(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockFraudDomainService struct {
	mock.Mock
}

func (m *MockFraudDomainService) AnalyzeFraud(ctx context.Context, userID *uint, email, ipAddress, userAgent string) (*entity.FraudAnalysis, error) {
	args := m.Called(ctx, userID, email, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if analysis, ok := args.Get(0).(*entity.FraudAnalysis); ok {
		return analysis, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFraudDomainService) CreateSecurityEvent(ctx context.Context, userID *uint, eventType, description, ipAddress, userAgent, severity string) error {
	args := m.Called(ctx, userID, eventType, description, ipAddress, userAgent, severity)
	return args.Error(0)
}

func (m *MockFraudDomainService) RecordLoginAttempt(ctx context.Context, email, ipAddress, userAgent string, success bool, failureReason string) error {
	args := m.Called(ctx, email, ipAddress, userAgent, success, failureReason)
	return args.Error(0)
}

func (m *MockFraudDomainService) DeactivateUserSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetFraudStats(ctx context.Context) (map[string]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if stats, ok := args.Get(0).(map[string]interface{}); ok {
		return stats, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFraudDomainService) AddIPToBlacklist(ctx context.Context, ip, reason, clientIP, userAgent string) error {
	args := m.Called(ctx, ip, reason, clientIP, userAgent)
	return args.Error(0)
}

func (m *MockFraudDomainService) RemoveIPFromBlacklist(ctx context.Context, ip, clientIP, userAgent string) error {
	args := m.Called(ctx, ip, clientIP, userAgent)
	return args.Error(0)
}

func (m *MockFraudDomainService) mockPaginatedResponse(methodName string, ctx context.Context, page, limit int) (interface{}, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) mockSimpleResponse(methodName string, ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	return args.Get(0), args.Error(1)
}

func (m *MockFraudDomainService) mockSessionOperation(methodName string, ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetBlacklistedIPs(ctx context.Context, page, limit int) (interface{}, error) {
	return m.mockPaginatedResponse("GetBlacklistedIPs", ctx, page, limit)
}

func (m *MockFraudDomainService) GetSecurityEvents(ctx context.Context, page, limit int) (interface{}, error) {
	return m.mockPaginatedResponse("GetSecurityEvents", ctx, page, limit)
}

func (m *MockFraudDomainService) CreateRateLimitRule(ctx context.Context, name, pattern string, maxRequests, windowSize int64) error {
	args := m.Called(ctx, name, pattern, maxRequests, windowSize)
	return args.Error(0)
}

func (m *MockFraudDomainService) UpdateRateLimitRule(ctx context.Context, id uint, name, pattern string, maxRequests, windowSize int64) error {
	args := m.Called(ctx, id, name, pattern, maxRequests, windowSize)
	return args.Error(0)
}

func (m *MockFraudDomainService) DeleteRateLimitRule(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockFraudDomainService) GetRateLimitRules(ctx context.Context) (interface{}, error) {
	return m.mockSimpleResponse("GetRateLimitRules", ctx)
}

func (m *MockFraudDomainService) GetActiveSessions(ctx context.Context) (interface{}, error) {
	return m.mockSimpleResponse("GetActiveSessions", ctx)
}

func (m *MockFraudDomainService) DeactivateSession(ctx context.Context, sessionID string) error {
	return m.mockSessionOperation("DeactivateSession", ctx, sessionID)
}

func (m *MockFraudDomainService) GetDevices(ctx context.Context) (interface{}, error) {
	return m.mockSimpleResponse("GetDevices", ctx)
}

func (m *MockFraudDomainService) TrustDevice(ctx context.Context, fingerprint string) error {
	args := m.Called(ctx, fingerprint)
	return args.Error(0)
}

func (m *MockFraudDomainService) CleanupExpiredData(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestAuthUsecaseRegister(t *testing.T) {
	tests := []struct {
		name      string
		req       usecase.RegisterRequest
		ipAddress string
		userAgent string
		setupMock func(*MockAuthDomainService, *MockFraudDomainService)
		wantErr   bool
	}{
		{
			name: "正常なユーザー登録",
			req: usecase.RegisterRequest{
				Name:     "テストユーザー",
				Email:    "test@example.com",
				Password: "password123",
				Age:      25,
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				user := entity.NewUser("テストユーザー", "test@example.com", 25)
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				roles := []string{"user"}
				fraudAnalysis := &entity.FraudAnalysis{RiskScore: 0.1}

				fraudService.On("AnalyzeFraud", ctx, (*uint)(nil), "test@example.com", "192.168.1.1", "test-agent").Return(fraudAnalysis, nil)
				authService.On("Register", ctx, "テストユーザー", "test@example.com", "password123", 25).Return(user, nil)
				fraudService.On("RecordLoginAttempt", ctx, "test@example.com", "192.168.1.1", "test-agent", true, "").Return(nil)
				fraudService.On("CreateSecurityEvent", ctx, &user.ID, "USER_REGISTRATION", "New user registered", "192.168.1.1", "test-agent", "LOW").Return(nil)
				authService.On("Login", ctx, "test@example.com", "password123").Return(auth, roles, nil)
				authService.On("GenerateAccessToken", uint(1), "test@example.com", roles).Return("access-token", nil)
				authService.On("GenerateRefreshToken", ctx, uint(1)).Return("refresh-token", nil)
			},
			wantErr: false,
		},
		{
			name: "高リスクユーザー登録でブロック",
			req: usecase.RegisterRequest{
				Name:     "テストユーザー",
				Email:    "suspicious@example.com",
				Password: "password123",
				Age:      25,
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				fraudAnalysis := entity.NewFraudAnalysis(0.9, []string{"suspicious IP", "unusual pattern"})

				fraudService.On("AnalyzeFraud", ctx, (*uint)(nil), "suspicious@example.com", "192.168.1.1", "test-agent").Return(fraudAnalysis, nil)
				fraudService.On("CreateSecurityEvent", ctx, (*uint)(nil), "HIGH_RISK_REGISTRATION", "High risk registration attempt", "192.168.1.1", "test-agent", "HIGH").Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService)

			ctx := context.Background()
			result, err := usecase.Register(ctx, tt.req, tt.ipAddress, tt.userAgent)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Email, result.User.Email)
				assert.Equal(t, tt.req.Name, result.User.Name)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			authService.AssertExpectations(t)
			fraudService.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseLogin(t *testing.T) {
	tests := []struct {
		name      string
		req       usecase.LoginRequest
		ipAddress string
		userAgent string
		setupMock func(*MockAuthDomainService, *MockFraudDomainService)
		wantErr   bool
	}{
		{
			name: "正常なログイン",
			req: usecase.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				roles := []string{"user"}
				fraudAnalysis := entity.NewFraudAnalysis(0.1, []string{"normal pattern"})

				fraudService.On("AnalyzeFraud", ctx, (*uint)(nil), "test@example.com", "192.168.1.1", "test-agent").Return(fraudAnalysis, nil)
				authService.On("Login", ctx, "test@example.com", "password123").Return(auth, roles, nil)
				fraudService.On("RecordLoginAttempt", ctx, "test@example.com", "192.168.1.1", "test-agent", true, "").Return(nil)
				fraudService.On("CreateSecurityEvent", ctx, &auth.UserID, "LOGIN", "User logged in successfully", "192.168.1.1", "test-agent", "LOW").Return(nil)
				authService.On("GenerateAccessToken", uint(1), "test@example.com", roles).Return("access-token", nil)
				authService.On("GenerateRefreshToken", ctx, uint(1)).Return("refresh-token", nil)
			},
			wantErr: false,
		},
		{
			name: "高リスクログインでブロック",
			req: usecase.LoginRequest{
				Email:    "suspicious@example.com",
				Password: "password123",
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				fraudAnalysis := entity.NewFraudAnalysis(0.9, []string{"suspicious IP", "unusual pattern"})

				fraudService.On("AnalyzeFraud", ctx, (*uint)(nil), "suspicious@example.com", "192.168.1.1", "test-agent").Return(fraudAnalysis, nil)
				fraudService.On("RecordLoginAttempt", ctx, "suspicious@example.com", "192.168.1.1", "test-agent", false, "High risk login blocked").Return(nil)
				fraudService.On("CreateSecurityEvent", ctx, (*uint)(nil), "HIGH_RISK_LOGIN", "High risk login attempt blocked", "192.168.1.1", "test-agent", "HIGH").Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService)

			ctx := context.Background()
			result, err := usecase.Login(ctx, tt.req, tt.ipAddress, tt.userAgent)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			authService.AssertExpectations(t)
			fraudService.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		req       usecase.RefreshTokenRequest
		setupMock func(*MockAuthDomainService, *MockFraudDomainService)
		wantErr   bool
	}{
		{
			name: "有効なリフレッシュトークン",
			req: usecase.RefreshTokenRequest{
				RefreshToken: "valid-refresh-token",
			},
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				auth, _ := entity.NewAuth(1, "test@example.com", "password123")
				roles := []string{"user"}
				newRefreshToken := "new-refresh-token"

				authService.On("RefreshToken", ctx, "valid-refresh-token").Return(auth, roles, newRefreshToken, nil)
				authService.On("GenerateAccessToken", uint(1), "test@example.com", roles).Return("new-access-token", nil)
			},
			wantErr: false,
		},
		{
			name: "無効なリフレッシュトークン",
			req: usecase.RefreshTokenRequest{
				RefreshToken: "invalid-refresh-token",
			},
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				authService.On("RefreshToken", ctx, "invalid-refresh-token").Return((*entity.Auth)(nil), []string(nil), "", assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService)

			ctx := context.Background()
			result, err := usecase.RefreshToken(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			authService.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseChangePassword(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		req       usecase.ChangePasswordRequest
		ipAddress string
		userAgent string
		setupMock func(*MockAuthDomainService, *MockFraudDomainService)
		wantErr   bool
	}{
		{
			name:   "正常なパスワード変更",
			userID: 1,
			req: usecase.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword123",
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				authService.On("ChangePassword", ctx, uint(1), "oldpassword", "newpassword123").Return(nil)
				fraudService.On("CreateSecurityEvent", ctx, &[]uint{1}[0], "PASSWORD_CHANGE", "User changed password", "192.168.1.1", "test-agent", "MEDIUM").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "現在のパスワードが間違っている",
			userID: 1,
			req: usecase.ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword123",
			},
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				authService.On("ChangePassword", ctx, uint(1), "wrongpassword", "newpassword123").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService)

			ctx := context.Background()
			err := usecase.ChangePassword(ctx, tt.userID, tt.req, tt.ipAddress, tt.userAgent)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authService.AssertExpectations(t)
			fraudService.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseLogout(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		sessionID string
		ipAddress string
		userAgent string
		setupMock func(*MockAuthDomainService, *MockFraudDomainService)
		wantErr   bool
	}{
		{
			name:      "正常なログアウト",
			userID:    1,
			sessionID: "session123",
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				authService.On("Logout", ctx, uint(1)).Return(nil)
				fraudService.On("DeactivateUserSession", ctx, "session123").Return(nil)
				fraudService.On("CreateSecurityEvent", ctx, &[]uint{1}[0], "LOGOUT", "User logged out", "192.168.1.1", "test-agent", "LOW").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "ログアウト時のエラー",
			userID:    1,
			sessionID: "session123",
			ipAddress: "192.168.1.1",
			userAgent: "test-agent",
			setupMock: func(authService *MockAuthDomainService, fraudService *MockFraudDomainService) {
				ctx := context.Background()
				authService.On("Logout", ctx, uint(1)).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService)

			ctx := context.Background()
			err := usecase.Logout(ctx, tt.userID, tt.sessionID, tt.ipAddress, tt.userAgent)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			authService.AssertExpectations(t)
			fraudService.AssertExpectations(t)
		})
	}
}
