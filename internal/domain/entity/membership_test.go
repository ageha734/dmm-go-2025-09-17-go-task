package entity_test

import (
	"testing"
	"time"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewMembershipTier(t *testing.T) {
	t.Run("新しいメンバーシップティアを正常に作成できる", func(t *testing.T) {
		name := "ゴールド"
		level := 2
		description := "ゴールドメンバー向けの特典"
		benefits := "送料無料、優先サポート"
		requirements := "年間10万円以上の購入"

		tier := entity.NewMembershipTier(name, level, description, benefits, requirements)

		assert.NotNil(t, tier)
		assert.Equal(t, name, tier.Name)
		assert.Equal(t, level, tier.Level)
		assert.Equal(t, description, tier.Description)
		assert.Equal(t, benefits, tier.Benefits)
		assert.Equal(t, requirements, tier.Requirements)
		assert.True(t, tier.IsActive)
		assert.False(t, tier.CreatedAt.IsZero())
		assert.False(t, tier.UpdatedAt.IsZero())
	})

	t.Run("空の値でもメンバーシップティアを作成できる", func(t *testing.T) {
		tier := entity.NewMembershipTier("", 0, "", "", "")

		assert.NotNil(t, tier)
		assert.Equal(t, "", tier.Name)
		assert.Equal(t, 0, tier.Level)
		assert.True(t, tier.IsActive)
	})
}

func TestNewUserMembership(t *testing.T) {
	t.Run("新しいユーザーメンバーシップを正常に作成できる", func(t *testing.T) {
		userID := uint(1)
		tierID := uint(2)

		membership := entity.NewUserMembership(userID, tierID)

		assert.NotNil(t, membership)
		assert.Equal(t, userID, membership.UserID)
		assert.Equal(t, tierID, membership.TierID)
		assert.Equal(t, 0, membership.Points)
		assert.Equal(t, 0.0, membership.TotalSpent)
		assert.False(t, membership.JoinedAt.IsZero())
		assert.Nil(t, membership.LastActivityAt)
		assert.Nil(t, membership.ExpiresAt)
		assert.True(t, membership.IsActive)
		assert.False(t, membership.CreatedAt.IsZero())
		assert.False(t, membership.UpdatedAt.IsZero())
	})
}

func TestUserMembershipAddPoints(t *testing.T) {
	membership := entity.NewUserMembership(1, 1)

	t.Run("正の値のポイントを正常に追加できる", func(t *testing.T) {
		initialPoints := membership.Points
		pointsToAdd := 100

		err := membership.AddPoints(pointsToAdd)

		assert.NoError(t, err)
		assert.Equal(t, initialPoints+pointsToAdd, membership.Points)
		assert.NotNil(t, membership.LastActivityAt)
	})

	t.Run("複数回ポイントを追加できる", func(t *testing.T) {
		initialPoints := membership.Points

		err1 := membership.AddPoints(50)
		err2 := membership.AddPoints(30)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, initialPoints+80, membership.Points)
	})

	t.Run("ゼロポイントの追加はエラーになる", func(t *testing.T) {
		initialPoints := membership.Points

		err := membership.AddPoints(0)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidPointAmount, err)
		assert.Equal(t, initialPoints, membership.Points)
	})

	t.Run("負のポイントの追加はエラーになる", func(t *testing.T) {
		initialPoints := membership.Points

		err := membership.AddPoints(-10)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidPointAmount, err)
		assert.Equal(t, initialPoints, membership.Points)
	})
}

func TestUserMembershipSpendPoints(t *testing.T) {
	membership := entity.NewUserMembership(1, 1)
	_ = membership.AddPoints(100)

	t.Run("十分なポイントがある場合は正常に消費できる", func(t *testing.T) {
		initialPoints := membership.Points
		pointsToSpend := 30

		err := membership.SpendPoints(pointsToSpend)

		assert.NoError(t, err)
		assert.Equal(t, initialPoints-pointsToSpend, membership.Points)
		assert.NotNil(t, membership.LastActivityAt)
	})

	t.Run("全ポイントを消費できる", func(t *testing.T) {
		currentPoints := membership.Points

		err := membership.SpendPoints(currentPoints)

		assert.NoError(t, err)
		assert.Equal(t, 0, membership.Points)
	})

	t.Run("不十分なポイントの場合はエラーになる", func(t *testing.T) {
		_ = membership.AddPoints(50)
		initialPoints := membership.Points

		err := membership.SpendPoints(initialPoints + 1)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInsufficientPoints, err)
		assert.Equal(t, initialPoints, membership.Points)
	})

	t.Run("ゼロポイントの消費はエラーになる", func(t *testing.T) {
		initialPoints := membership.Points

		err := membership.SpendPoints(0)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidPointAmount, err)
		assert.Equal(t, initialPoints, membership.Points)
	})

	t.Run("負のポイントの消費はエラーになる", func(t *testing.T) {
		initialPoints := membership.Points

		err := membership.SpendPoints(-10)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidPointAmount, err)
		assert.Equal(t, initialPoints, membership.Points)
	})
}

func TestUserMembershipUpdateTotalSpent(t *testing.T) {
	membership := entity.NewUserMembership(1, 1)

	t.Run("正の金額で総支出額を更新できる", func(t *testing.T) {
		amount := 1000.50

		membership.UpdateTotalSpent(amount)

		assert.Equal(t, amount, membership.TotalSpent)
		assert.NotNil(t, membership.LastActivityAt)
	})

	t.Run("複数回の支出額を累積できる", func(t *testing.T) {
		initialSpent := membership.TotalSpent
		amount1 := 500.25
		amount2 := 300.75

		membership.UpdateTotalSpent(amount1)
		membership.UpdateTotalSpent(amount2)

		assert.Equal(t, initialSpent+amount1+amount2, membership.TotalSpent)
	})

	t.Run("ゼロ金額では更新されない", func(t *testing.T) {
		initialSpent := membership.TotalSpent
		initialActivity := membership.LastActivityAt

		membership.UpdateTotalSpent(0)

		assert.Equal(t, initialSpent, membership.TotalSpent)
		assert.Equal(t, initialActivity, membership.LastActivityAt)
	})

	t.Run("負の金額では更新されない", func(t *testing.T) {
		initialSpent := membership.TotalSpent
		initialActivity := membership.LastActivityAt

		membership.UpdateTotalSpent(-100)

		assert.Equal(t, initialSpent, membership.TotalSpent)
		assert.Equal(t, initialActivity, membership.LastActivityAt)
	})
}

func TestUserMembershipUpdateLastActivity(t *testing.T) {
	membership := entity.NewUserMembership(1, 1)

	t.Run("最終活動時刻を更新できる", func(t *testing.T) {
		initialActivity := membership.LastActivityAt
		initialUpdatedAt := membership.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		membership.UpdateLastActivity()

		assert.NotEqual(t, initialActivity, membership.LastActivityAt)
		assert.NotNil(t, membership.LastActivityAt)
		assert.True(t, membership.UpdatedAt.After(initialUpdatedAt))
	})
}

func TestUserMembershipDeactivate(t *testing.T) {
	membership := entity.NewUserMembership(1, 1)

	t.Run("メンバーシップを非アクティブにできる", func(t *testing.T) {
		initialUpdatedAt := membership.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		membership.Deactivate()

		assert.False(t, membership.IsActive)
		assert.True(t, membership.UpdatedAt.After(initialUpdatedAt))
	})
}

func TestNewPointTransaction(t *testing.T) {
	t.Run("新しいポイント取引を正常に作成できる", func(t *testing.T) {
		userID := uint(1)
		transactionType := "earn"
		points := 100
		description := "購入による獲得"

		transaction := entity.NewPointTransaction(userID, transactionType, points, description)

		assert.NotNil(t, transaction)
		assert.Equal(t, userID, transaction.UserID)
		assert.Equal(t, transactionType, transaction.Type)
		assert.Equal(t, points, transaction.Points)
		assert.Equal(t, description, transaction.Description)
		assert.False(t, transaction.CreatedAt.IsZero())
	})

	t.Run("負のポイントでも取引を作成できる", func(t *testing.T) {
		transaction := entity.NewPointTransaction(1, "spend", -50, "ポイント使用")

		assert.NotNil(t, transaction)
		assert.Equal(t, -50, transaction.Points)
	})

	t.Run("空の説明でも取引を作成できる", func(t *testing.T) {
		transaction := entity.NewPointTransaction(1, "bonus", 25, "")

		assert.NotNil(t, transaction)
		assert.Equal(t, "", transaction.Description)
	})
}

func TestNewUserProfile(t *testing.T) {
	t.Run("新しいユーザープロフィールを正常に作成できる", func(t *testing.T) {
		userID := uint(1)

		profile := entity.NewUserProfile(userID)

		assert.NotNil(t, profile)
		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, "", profile.FirstName)
		assert.Equal(t, "", profile.LastName)
		assert.Equal(t, "", profile.PhoneNumber)
		assert.Nil(t, profile.DateOfBirth)
		assert.Equal(t, "", profile.Gender)
		assert.Nil(t, profile.Address)
		assert.Nil(t, profile.Preferences)
		assert.Equal(t, "", profile.Avatar)
		assert.Equal(t, "", profile.Bio)
		assert.False(t, profile.IsVerified)
		assert.Nil(t, profile.VerifiedAt)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())
	})
}

func TestUserProfileUpdateProfile(t *testing.T) {
	profile := entity.NewUserProfile(1)

	t.Run("プロフィール情報を正常に更新できる", func(t *testing.T) {
		firstName := "太郎"
		lastName := "田中"
		phoneNumber := "090-1234-5678"
		gender := "male"
		bio := "よろしくお願いします"
		dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		initialUpdatedAt := profile.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		profile.UpdateProfile(firstName, lastName, phoneNumber, gender, bio, &dateOfBirth)

		assert.Equal(t, firstName, profile.FirstName)
		assert.Equal(t, lastName, profile.LastName)
		assert.Equal(t, phoneNumber, profile.PhoneNumber)
		assert.Equal(t, gender, profile.Gender)
		assert.Equal(t, bio, profile.Bio)
		assert.Equal(t, dateOfBirth, *profile.DateOfBirth)
		assert.True(t, profile.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("空の値では既存の値が保持される", func(t *testing.T) {
		profile.FirstName = "既存の名前"
		profile.LastName = "既存の姓"

		profile.UpdateProfile("", "", "", "", "", nil)

		assert.Equal(t, "既存の名前", profile.FirstName)
		assert.Equal(t, "既存の姓", profile.LastName)
	})

	t.Run("一部の値のみ更新できる", func(t *testing.T) {
		originalFirstName := profile.FirstName
		newBio := "新しい自己紹介"

		profile.UpdateProfile("", "", "", "", newBio, nil)

		assert.Equal(t, originalFirstName, profile.FirstName)
		assert.Equal(t, newBio, profile.Bio)
	})
}

func TestUserProfileVerify(t *testing.T) {
	profile := entity.NewUserProfile(1)

	t.Run("プロフィールを認証済みにできる", func(t *testing.T) {
		initialUpdatedAt := profile.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		profile.Verify()

		assert.True(t, profile.IsVerified)
		assert.NotNil(t, profile.VerifiedAt)
		assert.False(t, profile.VerifiedAt.IsZero())
		assert.True(t, profile.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("既に認証済みのプロフィールも再認証できる", func(t *testing.T) {
		firstVerifiedAt := *profile.VerifiedAt

		time.Sleep(1 * time.Millisecond)
		profile.Verify()

		assert.True(t, profile.IsVerified)
		assert.True(t, profile.VerifiedAt.After(firstVerifiedAt))
	})
}

func TestNewNotification(t *testing.T) {
	t.Run("新しい通知を正常に作成できる", func(t *testing.T) {
		userID := uint(1)
		notificationType := "info"
		title := "お知らせ"
		message := "新機能がリリースされました"
		data := `{"feature": "new_dashboard"}`

		notification := entity.NewNotification(userID, notificationType, title, message, &data)

		assert.NotNil(t, notification)
		assert.Equal(t, userID, notification.UserID)
		assert.Equal(t, notificationType, notification.Type)
		assert.Equal(t, title, notification.Title)
		assert.Equal(t, message, notification.Message)
		assert.Equal(t, data, *notification.Data)
		assert.False(t, notification.IsRead)
		assert.Nil(t, notification.ReadAt)
		assert.False(t, notification.CreatedAt.IsZero())
		assert.False(t, notification.UpdatedAt.IsZero())
	})

	t.Run("データなしで通知を作成できる", func(t *testing.T) {
		notification := entity.NewNotification(1, "warning", "警告", "注意が必要です", nil)

		assert.NotNil(t, notification)
		assert.Nil(t, notification.Data)
	})

	t.Run("空のタイトルとメッセージでも通知を作成できる", func(t *testing.T) {
		notification := entity.NewNotification(1, "system", "", "", nil)

		assert.NotNil(t, notification)
		assert.Equal(t, "", notification.Title)
		assert.Equal(t, "", notification.Message)
	})
}

func TestNotificationMarkAsRead(t *testing.T) {
	notification := entity.NewNotification(1, "info", "テスト", "テストメッセージ", nil)

	t.Run("未読の通知を既読にできる", func(t *testing.T) {
		initialUpdatedAt := notification.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		notification.MarkAsRead()

		assert.True(t, notification.IsRead)
		assert.NotNil(t, notification.ReadAt)
		assert.False(t, notification.ReadAt.IsZero())
		assert.True(t, notification.UpdatedAt.After(initialUpdatedAt))
	})

	t.Run("既読の通知を再度既読にしても変更されない", func(t *testing.T) {
		firstReadAt := *notification.ReadAt
		firstUpdatedAt := notification.UpdatedAt

		time.Sleep(1 * time.Millisecond)
		notification.MarkAsRead()

		assert.True(t, notification.IsRead)
		assert.Equal(t, firstReadAt, *notification.ReadAt)
		assert.Equal(t, firstUpdatedAt, notification.UpdatedAt)
	})
}
