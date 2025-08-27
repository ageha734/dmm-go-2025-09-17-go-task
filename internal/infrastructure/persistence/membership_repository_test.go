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

func setupMembershipRepositoryTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestUserMembershipRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	membership := &entity.UserMembership{
		UserID:   1,
		TierID:   1,
		Points:   100,
		IsActive: true,
		JoinedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `user_memberships`").
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
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, membership)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), membership.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserMembershipRepositoryGetByUserID(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)
	expectedMembership := &entity.UserMembership{
		ID:        1,
		UserID:    userID,
		TierID:    1,
		Points:    100,
		IsActive:  true,
		JoinedAt:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "tier_id", "points", "total_spent",
		"joined_at", "last_activity_at", "expires_at", "is_active",
		"created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedMembership.ID,
		expectedMembership.UserID,
		expectedMembership.TierID,
		expectedMembership.Points,
		0.0,
		expectedMembership.JoinedAt,
		nil,
		nil,
		expectedMembership.IsActive,
		expectedMembership.CreatedAt,
		expectedMembership.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `user_memberships` WHERE user_id = \\? AND `user_memberships`.`deleted_at` IS NULL ORDER BY `user_memberships`.`id` LIMIT \\?").
		WithArgs(userID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedMembership.ID, result.ID)
	assert.Equal(t, expectedMembership.UserID, result.UserID)
	assert.Equal(t, expectedMembership.Points, result.Points)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserMembershipRepositoryUpdate(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	membership := &entity.UserMembership{
		ID:        1,
		UserID:    1,
		TierID:    2,
		Points:    200,
		IsActive:  true,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_memberships` SET").
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
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(ctx, membership)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserMembershipRepositoryDelete(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_memberships` SET `deleted_at`=\\? WHERE user_id = \\? AND `user_memberships`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserMembershipRepositoryGetStats(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(150)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `user_memberships` WHERE `user_memberships`.`deleted_at` IS NULL").
		WillReturnRows(rows)

	stats, err := repo.GetStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(150), stats["total_members"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserMembershipRepositoryList(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserMembershipRepository(gormDB)
	ctx := context.Background()

	offset := 0
	limit := 10

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `user_memberships` WHERE `user_memberships`.`deleted_at` IS NULL").
		WillReturnRows(countRows)

	memberships := []entity.UserMembership{
		{
			ID:        1,
			UserID:    1,
			TierID:    1,
			Points:    100,
			IsActive:  true,
			JoinedAt:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UserID:    2,
			TierID:    2,
			Points:    200,
			IsActive:  true,
			JoinedAt:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "tier_id", "points", "total_spent",
		"joined_at", "last_activity_at", "expires_at", "is_active",
		"created_at", "updated_at", "deleted_at",
	})
	for _, membership := range memberships {
		rows.AddRow(
			membership.ID,
			membership.UserID,
			membership.TierID,
			membership.Points,
			0.0,
			membership.JoinedAt,
			nil,
			nil,
			membership.IsActive,
			membership.CreatedAt,
			membership.UpdatedAt,
			nil,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM `user_memberships` WHERE `user_memberships`.`deleted_at` IS NULL LIMIT \\?").
		WithArgs(limit).
		WillReturnRows(rows)

	result, total, err := repo.List(ctx, offset, limit)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	assert.Equal(t, memberships[0].ID, result[0].ID)
	assert.Equal(t, memberships[1].ID, result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserProfileRepository(gormDB)
	ctx := context.Background()

	profile := &entity.UserProfile{
		UserID:     1,
		FirstName:  "John",
		LastName:   "Doe",
		IsVerified: false,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `user_profiles`").
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
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, profile)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), profile.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepositoryGetByUserID(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserProfileRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)
	expectedProfile := &entity.UserProfile{
		ID:         1,
		UserID:     userID,
		FirstName:  "John",
		LastName:   "Doe",
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "first_name", "last_name", "phone_number",
		"date_of_birth", "gender", "address", "preferences", "avatar",
		"bio", "is_verified", "verified_at", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedProfile.ID,
		expectedProfile.UserID,
		expectedProfile.FirstName,
		expectedProfile.LastName,
		"",
		nil,
		"",
		nil,
		nil,
		"",
		"",
		expectedProfile.IsVerified,
		nil,
		expectedProfile.CreatedAt,
		expectedProfile.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `user_profiles` WHERE user_id = \\? AND `user_profiles`.`deleted_at` IS NULL ORDER BY `user_profiles`.`id` LIMIT \\?").
		WithArgs(userID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByUserID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProfile.ID, result.ID)
	assert.Equal(t, expectedProfile.UserID, result.UserID)
	assert.Equal(t, expectedProfile.FirstName, result.FirstName)
	assert.Equal(t, expectedProfile.LastName, result.LastName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepositoryUpdate(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserProfileRepository(gormDB)
	ctx := context.Background()

	profile := &entity.UserProfile{
		ID:         1,
		UserID:     1,
		FirstName:  "Jane",
		LastName:   "Smith",
		IsVerified: true,
		UpdatedAt:  time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_profiles` SET").
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

	err := repo.Update(ctx, profile)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserProfileRepositoryDelete(t *testing.T) {
	gormDB, mock, cleanup := setupMembershipRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserProfileRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_profiles` SET `deleted_at`=\\? WHERE user_id = \\? AND `user_profiles`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
