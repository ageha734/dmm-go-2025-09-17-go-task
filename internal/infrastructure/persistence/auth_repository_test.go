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

func setupAuthRepositoryTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestAuthRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	auth := &entity.Auth{
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsActive:     true,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `auths`").
		WithArgs(
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

	err := repo.Create(ctx, auth)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), auth.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepositoryGetByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"
	expectedAuth := &entity.Auth{
		ID:           1,
		UserID:       1,
		Email:        email,
		PasswordHash: "hashedpassword",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "email", "password_hash", "is_active",
		"last_login_at", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedAuth.ID,
		expectedAuth.UserID,
		expectedAuth.Email,
		expectedAuth.PasswordHash,
		expectedAuth.IsActive,
		nil,
		expectedAuth.CreatedAt,
		expectedAuth.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `auths` WHERE email = \\? AND `auths`.`deleted_at` IS NULL ORDER BY `auths`.`id` LIMIT \\?").
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.GetByEmail(ctx, email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAuth.ID, result.ID)
	assert.Equal(t, expectedAuth.Email, result.Email)
	assert.Equal(t, expectedAuth.UserID, result.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepositoryGetByUserID(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)
	expectedAuth := &entity.Auth{
		ID:           1,
		UserID:       userID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "email", "password_hash", "is_active",
		"last_login_at", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedAuth.ID,
		expectedAuth.UserID,
		expectedAuth.Email,
		expectedAuth.PasswordHash,
		expectedAuth.IsActive,
		nil,
		expectedAuth.CreatedAt,
		expectedAuth.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `auths` WHERE user_id = \\? AND `auths`.`deleted_at` IS NULL ORDER BY `auths`.`id` LIMIT \\?").
		WithArgs(userID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAuth.ID, result.ID)
	assert.Equal(t, expectedAuth.UserID, result.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepositoryUpdate(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	auth := &entity.Auth{
		ID:           1,
		UserID:       1,
		Email:        "test@example.com",
		PasswordHash: "newhashedpassword",
		IsActive:     false,
		UpdatedAt:    time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auths` SET").
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

	err := repo.Update(ctx, auth)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepositoryDelete(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `auths` SET `deleted_at`=\\? WHERE user_id = \\? AND `auths`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthRepositoryExistsByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewAuthRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"

	t.Run("exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `auths` WHERE email = \\? AND `auths`.`deleted_at` IS NULL").
			WithArgs(email).
			WillReturnRows(rows)

		exists, err := repo.ExistsByEmail(ctx, email)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `auths` WHERE email = \\? AND `auths`.`deleted_at` IS NULL").
			WithArgs(email).
			WillReturnRows(rows)

		exists, err := repo.ExistsByEmail(ctx, email)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRefreshTokenRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewRefreshTokenRepository(gormDB)
	ctx := context.Background()

	token := &entity.RefreshToken{
		UserID:    1,
		Token:     "refresh_token_123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `refresh_tokens`").
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

	err := repo.Create(ctx, token)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), token.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRefreshTokenRepositoryGetByToken(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewRefreshTokenRepository(gormDB)
	ctx := context.Background()

	tokenStr := "refresh_token_123"
	expectedToken := &entity.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "token", "expires_at", "is_revoked",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedToken.ID,
		expectedToken.UserID,
		expectedToken.Token,
		expectedToken.ExpiresAt,
		expectedToken.IsRevoked,
		expectedToken.CreatedAt,
		expectedToken.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `refresh_tokens` WHERE token = \\? AND `refresh_tokens`.`deleted_at` IS NULL ORDER BY `refresh_tokens`.`id` LIMIT \\?").
		WithArgs(tokenStr, 1).
		WillReturnRows(rows)

	result, err := repo.GetByToken(ctx, tokenStr)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedToken.ID, result.ID)
	assert.Equal(t, expectedToken.Token, result.Token)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRefreshTokenRepositoryRevokeByUserID(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewRefreshTokenRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `refresh_tokens` SET `is_revoked`=\\?,`updated_at`=\\? WHERE user_id = \\? AND `refresh_tokens`.`deleted_at` IS NULL").
		WithArgs(true, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.RevokeByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRefreshTokenRepositoryDeleteExpired(t *testing.T) {
	gormDB, mock, cleanup := setupAuthRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewRefreshTokenRepository(gormDB)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `refresh_tokens` SET `deleted_at`=\\? WHERE expires_at < NOW\\(\\) AND `refresh_tokens`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteExpired(ctx)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
