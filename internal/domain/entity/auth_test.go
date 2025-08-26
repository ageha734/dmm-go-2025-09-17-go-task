package entity_test

import (
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAuth(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "正常な認証情報作成",
			userID:   1,
			email:    "test@example.com",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "弱いパスワードでエラー",
			userID:   1,
			email:    "test@example.com",
			password: "123",
			wantErr:  entity.ErrWeakPassword,
		},
		{
			name:     "空のパスワードでエラー",
			userID:   1,
			email:    "test@example.com",
			password: "",
			wantErr:  entity.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := entity.NewAuth(tt.userID, tt.email, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, auth)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, auth)
				assert.Equal(t, tt.userID, auth.UserID)
				assert.Equal(t, tt.email, auth.Email)
				assert.True(t, auth.IsActive)
				assert.NotEmpty(t, auth.PasswordHash)
				assert.False(t, auth.CreatedAt.IsZero())
				assert.False(t, auth.UpdatedAt.IsZero())
				assert.Nil(t, auth.LastLoginAt)

				err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(tt.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthVerifyPassword(t *testing.T) {
	auth, err := entity.NewAuth(1, "test@example.com", "password123")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "正しいパスワード",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "間違ったパスワード",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "空のパスワード",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.VerifyPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthChangePassword(t *testing.T) {
	tests := []struct {
		name            string
		currentPassword string
		newPassword     string
		wantErr         error
	}{
		{
			name:            "正常なパスワード変更",
			currentPassword: "oldpassword",
			newPassword:     "newpassword123",
			wantErr:         nil,
		},
		{
			name:            "現在のパスワードが間違っている",
			currentPassword: "wrongpassword",
			newPassword:     "newpassword123",
			wantErr:         entity.ErrInvalidPassword,
		},
		{
			name:            "新しいパスワードが弱い",
			currentPassword: "oldpassword",
			newPassword:     "123",
			wantErr:         entity.ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAuth, _ := entity.NewAuth(1, "test@example.com", "oldpassword")
			oldUpdatedAt := testAuth.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			err := testAuth.ChangePassword(tt.currentPassword, tt.newPassword)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, testAuth.UpdatedAt.After(oldUpdatedAt))

				err = testAuth.VerifyPassword(tt.newPassword)
				assert.NoError(t, err)

				err = testAuth.VerifyPassword(tt.currentPassword)
				assert.Error(t, err)
			}
		})
	}
}

func TestAuthUpdateLastLogin(t *testing.T) {
	auth, err := entity.NewAuth(1, "test@example.com", "password123")
	assert.NoError(t, err)

	assert.Nil(t, auth.LastLoginAt)

	oldUpdatedAt := auth.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	auth.UpdateLastLogin()

	assert.NotNil(t, auth.LastLoginAt)
	assert.True(t, auth.UpdatedAt.After(oldUpdatedAt))
	assert.True(t, auth.LastLoginAt.After(oldUpdatedAt))
}

func TestAuthDeactivate(t *testing.T) {
	auth, err := entity.NewAuth(1, "test@example.com", "password123")
	assert.NoError(t, err)

	assert.True(t, auth.IsActive)

	oldUpdatedAt := auth.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	auth.Deactivate()

	assert.False(t, auth.IsActive)
	assert.True(t, auth.UpdatedAt.After(oldUpdatedAt))
}

func TestAuthActivate(t *testing.T) {
	auth, err := entity.NewAuth(1, "test@example.com", "password123")
	assert.NoError(t, err)

	auth.Deactivate()
	assert.False(t, auth.IsActive)

	oldUpdatedAt := auth.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	auth.Activate()

	assert.True(t, auth.IsActive)
	assert.True(t, auth.UpdatedAt.After(oldUpdatedAt))
}

func TestNewRole(t *testing.T) {
	tests := []struct {
		name        string
		roleName    string
		description string
	}{
		{
			name:        "管理者ロール作成",
			roleName:    "admin",
			description: "管理者権限",
		},
		{
			name:        "ユーザーロール作成",
			roleName:    "user",
			description: "一般ユーザー権限",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role := entity.NewRole(tt.roleName, tt.description)

			assert.NotNil(t, role)
			assert.Equal(t, tt.roleName, role.Name)
			assert.Equal(t, tt.description, role.Description)
			assert.False(t, role.CreatedAt.IsZero())
			assert.False(t, role.UpdatedAt.IsZero())
		})
	}
}

func TestNewRefreshToken(t *testing.T) {
	userID := uint(1)
	token := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	refreshToken := entity.NewRefreshToken(userID, token, expiresAt)

	assert.NotNil(t, refreshToken)
	assert.Equal(t, userID, refreshToken.UserID)
	assert.Equal(t, token, refreshToken.Token)
	assert.Equal(t, expiresAt, refreshToken.ExpiresAt)
	assert.False(t, refreshToken.IsRevoked)
	assert.False(t, refreshToken.CreatedAt.IsZero())
	assert.False(t, refreshToken.UpdatedAt.IsZero())
}

func TestRefreshTokenRevoke(t *testing.T) {
	refreshToken := entity.NewRefreshToken(1, "test-token", time.Now().Add(24*time.Hour))

	assert.False(t, refreshToken.IsRevoked)

	oldUpdatedAt := refreshToken.UpdatedAt
	time.Sleep(1 * time.Millisecond)

	refreshToken.Revoke()

	assert.True(t, refreshToken.IsRevoked)
	assert.True(t, refreshToken.UpdatedAt.After(oldUpdatedAt))
}

func TestRefreshTokenIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		want      bool
	}{
		{
			name:      "未来の時刻で有効",
			expiresAt: time.Now().Add(1 * time.Hour),
			want:      false,
		},
		{
			name:      "過去の時刻で期限切れ",
			expiresAt: time.Now().Add(-1 * time.Hour),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshToken := entity.NewRefreshToken(1, "test-token", tt.expiresAt)
			got := refreshToken.IsExpired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRefreshTokenIsValid(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		isRevoked bool
		want      bool
	}{
		{
			name:      "有効なトークン",
			expiresAt: time.Now().Add(1 * time.Hour),
			isRevoked: false,
			want:      true,
		},
		{
			name:      "無効化されたトークン",
			expiresAt: time.Now().Add(1 * time.Hour),
			isRevoked: true,
			want:      false,
		},
		{
			name:      "期限切れのトークン",
			expiresAt: time.Now().Add(-1 * time.Hour),
			isRevoked: false,
			want:      false,
		},
		{
			name:      "無効化され期限切れのトークン",
			expiresAt: time.Now().Add(-1 * time.Hour),
			isRevoked: true,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshToken := entity.NewRefreshToken(1, "test-token", tt.expiresAt)
			if tt.isRevoked {
				refreshToken.Revoke()
			}
			got := refreshToken.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}
