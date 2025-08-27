package usecase_test

import (
	"context"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/usecase"
	"github.com/stretchr/testify/assert"
)

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

			usecase := usecase.NewAuthUsecase(authService, fraudService, nil)

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

			usecase := usecase.NewAuthUsecase(authService, fraudService, nil)

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

			usecase := usecase.NewAuthUsecase(authService, fraudService, nil)

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

			usecase := usecase.NewAuthUsecase(authService, fraudService, nil)

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
				authService.On("Logout", ctx, uint(1), "test-token").Return(nil)
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
				authService.On("Logout", ctx, uint(1), "test-token").Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := new(MockAuthDomainService)
			fraudService := new(MockFraudDomainService)
			tt.setupMock(authService, fraudService)

			usecase := usecase.NewAuthUsecase(authService, fraudService, nil)

			ctx := context.Background()
			err := usecase.Logout(ctx, tt.userID, "test-token", tt.sessionID, tt.ipAddress, tt.userAgent)

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
