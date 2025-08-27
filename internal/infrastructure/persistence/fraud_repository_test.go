package persistence_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/domain/entity"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupFraudRepositoryTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	cleanup := func() {
		_ = db.Close()
	}

	return gormDB, mock, cleanup
}

func TestSecurityEventRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewSecurityEventRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)
	event := &entity.SecurityEvent{
		UserID:      &userID,
		EventType:   "login_failed",
		Description: "Failed login attempt",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Mozilla/5.0",
		Severity:    "medium",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `security_events`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, event)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), event.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSecurityEventRepositoryGetByID(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewSecurityEventRepository(gormDB)
	ctx := context.Background()

	eventID := uint(1)
	userID := uint(1)
	expectedEvent := &entity.SecurityEvent{
		ID:          eventID,
		UserID:      &userID,
		EventType:   "login_failed",
		Description: "Failed login attempt",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Mozilla/5.0",
		Severity:    "medium",
		CreatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "event_type", "description", "ip_address",
		"user_agent", "severity", "metadata", "created_at", "deleted_at",
	}).AddRow(
		expectedEvent.ID,
		expectedEvent.UserID,
		expectedEvent.EventType,
		expectedEvent.Description,
		expectedEvent.IPAddress,
		expectedEvent.UserAgent,
		expectedEvent.Severity,
		nil,
		expectedEvent.CreatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `security_events` WHERE `security_events`.`id` = \\? AND `security_events`.`deleted_at` IS NULL ORDER BY `security_events`.`id` LIMIT \\?").
		WithArgs(eventID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, eventID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEvent.ID, result.ID)
	assert.Equal(t, expectedEvent.EventType, result.EventType)
	assert.Equal(t, expectedEvent.IPAddress, result.IPAddress)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIPBlacklistRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewIPBlacklistRepository(gormDB)
	ctx := context.Background()

	blacklist := &entity.IPBlacklist{
		IPAddress: "192.168.1.100",
		Reason:    "Suspicious activity",
		IsActive:  true,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `ip_blacklists`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, blacklist)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), blacklist.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIPBlacklistRepositoryGetByIP(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewIPBlacklistRepository(gormDB)
	ctx := context.Background()

	ipAddress := "192.168.1.100"
	expectedBlacklist := &entity.IPBlacklist{
		ID:        1,
		IPAddress: ipAddress,
		Reason:    "Suspicious activity",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "ip_address", "reason", "expires_at", "is_active",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedBlacklist.ID,
		expectedBlacklist.IPAddress,
		expectedBlacklist.Reason,
		nil,
		expectedBlacklist.IsActive,
		expectedBlacklist.CreatedAt,
		expectedBlacklist.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `ip_blacklists` WHERE ip_address = \\? AND `ip_blacklists`.`deleted_at` IS NULL ORDER BY `ip_blacklists`.`id` LIMIT \\?").
		WithArgs(ipAddress, 1).
		WillReturnRows(rows)

	result, err := repo.GetByIP(ctx, ipAddress)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBlacklist.ID, result.ID)
	assert.Equal(t, expectedBlacklist.IPAddress, result.IPAddress)
	assert.Equal(t, expectedBlacklist.Reason, result.Reason)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIPBlacklistRepositoryIsBlacklisted(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewIPBlacklistRepository(gormDB)
	ctx := context.Background()

	ipAddress := "192.168.1.100"

	t.Run("is blacklisted", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `ip_blacklists` WHERE \\(ip_address = \\? AND is_active = \\? AND \\(expires_at IS NULL OR expires_at > NOW\\(\\)\\)\\) AND `ip_blacklists`.`deleted_at` IS NULL").
			WithArgs(ipAddress, true).
			WillReturnRows(rows)

		isBlacklisted, err := repo.IsBlacklisted(ctx, ipAddress)

		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not blacklisted", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `ip_blacklists` WHERE \\(ip_address = \\? AND is_active = \\? AND \\(expires_at IS NULL OR expires_at > NOW\\(\\)\\)\\) AND `ip_blacklists`.`deleted_at` IS NULL").
			WithArgs(ipAddress, true).
			WillReturnRows(rows)

		isBlacklisted, err := repo.IsBlacklisted(ctx, ipAddress)

		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIPBlacklistRepositoryCleanupExpired(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewIPBlacklistRepository(gormDB)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `ip_blacklists` SET `is_active`=\\?,`updated_at`=\\? WHERE \\(expires_at IS NOT NULL AND expires_at <= NOW\\(\\)\\) AND `ip_blacklists`.`deleted_at` IS NULL").
		WithArgs(false, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.CleanupExpired(ctx)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewLoginAttemptRepository(gormDB)
	ctx := context.Background()

	attempt := &entity.LoginAttempt{
		Email:      "test@example.com",
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0",
		Success:    false,
		FailReason: "Invalid password",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `login_attempts`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, attempt)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), attempt.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptRepositoryGetByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewLoginAttemptRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"
	since := time.Now().Add(-1 * time.Hour)

	attempts := []entity.LoginAttempt{
		{
			ID:         1,
			Email:      email,
			IPAddress:  "192.168.1.1",
			UserAgent:  "Mozilla/5.0",
			Success:    false,
			FailReason: "Invalid password",
			CreatedAt:  time.Now(),
		},
		{
			ID:         2,
			Email:      email,
			IPAddress:  "192.168.1.2",
			UserAgent:  "Chrome/90.0",
			Success:    true,
			FailReason: "",
			CreatedAt:  time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "email", "ip_address", "user_agent", "success",
		"fail_reason", "created_at", "deleted_at",
	})
	for _, attempt := range attempts {
		rows.AddRow(
			attempt.ID,
			attempt.Email,
			attempt.IPAddress,
			attempt.UserAgent,
			attempt.Success,
			attempt.FailReason,
			attempt.CreatedAt,
			nil,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM `login_attempts` WHERE \\(email = \\? AND created_at >= \\?\\) AND `login_attempts`.`deleted_at` IS NULL ORDER BY created_at DESC").
		WithArgs(email, since).
		WillReturnRows(rows)

	result, err := repo.GetByEmail(ctx, email, since)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, attempts[0].ID, result[0].ID)
	assert.Equal(t, attempts[1].ID, result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptRepositoryCountFailedAttempts(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewLoginAttemptRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"
	since := time.Now().Add(-1 * time.Hour)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `login_attempts` WHERE \\(email = \\? AND success = \\? AND created_at >= \\?\\) AND `login_attempts`.`deleted_at` IS NULL").
		WithArgs(email, false, since).
		WillReturnRows(rows)

	count, err := repo.CountFailedAttempts(ctx, email, since)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLoginAttemptRepositoryCleanupOld(t *testing.T) {
	gormDB, mock, cleanup := setupFraudRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewLoginAttemptRepository(gormDB)
	ctx := context.Background()

	before := time.Now().Add(-30 * 24 * time.Hour)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `login_attempts` SET `deleted_at`=\\? WHERE created_at < \\? AND `login_attempts`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), before).
		WillReturnResult(sqlmock.NewResult(1, 10))
	mock.ExpectCommit()

	err := repo.CleanupOld(ctx, before)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
