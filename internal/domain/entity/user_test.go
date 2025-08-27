package entity_test

import (
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name  string
		uname string
		email string
		age   int
	}{
		{
			name:  "正常なユーザー作成",
			uname: "テストユーザー",
			email: "test@example.com",
			age:   25,
		},
		{
			name:  "空の名前でもユーザー作成",
			uname: "",
			email: "test@example.com",
			age:   30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := entity.NewUser(tt.uname, tt.email, tt.age)

			assert.NotNil(t, user)
			assert.Equal(t, tt.uname, user.Name)
			assert.Equal(t, tt.email, user.Email)
			assert.Equal(t, tt.age, user.Age)
			assert.False(t, user.CreatedAt.IsZero())
			assert.False(t, user.UpdatedAt.IsZero())
		})
	}
}

func TestUserUpdateProfile(t *testing.T) {
	tests := []struct {
		name     string
		user     *entity.User
		newName  string
		newAge   int
		wantName string
		wantAge  int
	}{
		{
			name: "正常なプロフィール更新",
			user: &entity.User{
				ID:    1,
				Name:  "旧名前",
				Email: "old@example.com",
				Age:   25,
			},
			newName:  "新名前",
			newAge:   30,
			wantName: "新名前",
			wantAge:  30,
		},
		{
			name: "空の名前は更新されない",
			user: &entity.User{
				ID:    1,
				Name:  "旧名前",
				Email: "old@example.com",
				Age:   25,
			},
			newName:  "",
			newAge:   30,
			wantName: "旧名前",
			wantAge:  30,
		},
		{
			name: "0以下の年齢は更新されない",
			user: &entity.User{
				ID:    1,
				Name:  "旧名前",
				Email: "old@example.com",
				Age:   25,
			},
			newName:  "新名前",
			newAge:   -5,
			wantName: "新名前",
			wantAge:  25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := tt.user.UpdatedAt
			time.Sleep(1 * time.Millisecond)

			tt.user.UpdateProfile(tt.newName, tt.newAge)

			assert.Equal(t, tt.wantName, tt.user.Name)
			assert.Equal(t, tt.wantAge, tt.user.Age)
			assert.True(t, tt.user.UpdatedAt.After(oldUpdatedAt))
		})
	}
}

func TestUserIsValidAge(t *testing.T) {
	tests := []struct {
		name string
		user *entity.User
		want bool
	}{
		{
			name: "有効な年齢（25歳）",
			user: &entity.User{
				Age: 25,
			},
			want: true,
		},
		{
			name: "有効な年齢（0歳）",
			user: &entity.User{
				Age: 0,
			},
			want: true,
		},
		{
			name: "有効な年齢（150歳）",
			user: &entity.User{
				Age: 150,
			},
			want: true,
		},
		{
			name: "無効な年齢（負の値）",
			user: &entity.User{
				Age: -1,
			},
			want: false,
		},
		{
			name: "無効な年齢（151歳）",
			user: &entity.User{
				Age: 151,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsValidAge()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserIsValidEmail(t *testing.T) {
	tests := []struct {
		name string
		user *entity.User
		want bool
	}{
		{
			name: "有効なメールアドレス",
			user: &entity.User{
				Email: "test@example.com",
			},
			want: true,
		},
		{
			name: "空のメールアドレス",
			user: &entity.User{
				Email: "",
			},
			want: false,
		},
		{
			name: "単一文字のメールアドレス",
			user: &entity.User{
				Email: "a",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.user.IsValidEmail()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserStructComparison(t *testing.T) {
	user1 := &entity.User{
		ID:        1,
		Name:      "テストユーザー",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	user2 := &entity.User{
		ID:        1,
		Name:      "テストユーザー",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if diff := cmp.Diff(user1, user2); diff != "" {
		t.Errorf("User mismatch (-want +got):\n%s", diff)
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(entity.User{}, "CreatedAt", "UpdatedAt"),
	}

	user3 := &entity.User{
		ID:        1,
		Name:      "テストユーザー",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if diff := cmp.Diff(user1, user3, opts); diff != "" {
		t.Errorf("User mismatch (ignoring timestamps) (-want +got):\n%s", diff)
	}
}
