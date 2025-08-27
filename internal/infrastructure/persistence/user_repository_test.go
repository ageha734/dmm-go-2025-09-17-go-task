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

func setupUserRepositoryTest(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestUserRepositoryCreate(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	user := &entity.User{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByID(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)
	expectedUser := &entity.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Age:       25,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "age", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedUser.ID,
		expectedUser.Name,
		expectedUser.Email,
		expectedUser.Age,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(userID, 1).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Name, result.Name)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.Age, result.Age)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryGetByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"
	expectedUser := &entity.User{
		ID:        1,
		Name:      "Test User",
		Email:     email,
		Age:       25,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "age", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		expectedUser.ID,
		expectedUser.Name,
		expectedUser.Email,
		expectedUser.Age,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
		nil,
	)

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE email = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?").
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.GetByEmail(ctx, email)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryUpdate(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	user := &entity.User{
		ID:        1,
		Name:      "Updated User",
		Email:     "updated@example.com",
		Age:       30,
		UpdatedAt: time.Now(),
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET").
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

	err := repo.Update(ctx, user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryDelete(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users` SET `deleted_at`=\\? WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryList(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	offset := 0
	limit := 10

	countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE `users`.`deleted_at` IS NULL").
		WillReturnRows(countRows)

	users := []entity.User{
		{
			ID:        1,
			Name:      "User 1",
			Email:     "user1@example.com",
			Age:       25,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "User 2",
			Email:     "user2@example.com",
			Age:       30,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "name", "email", "age", "created_at", "updated_at", "deleted_at",
	})
	for _, user := range users {
		rows.AddRow(
			user.ID,
			user.Name,
			user.Email,
			user.Age,
			user.CreatedAt,
			user.UpdatedAt,
			nil,
		)
	}

	mock.ExpectQuery("SELECT \\* FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT \\?").
		WithArgs(limit).
		WillReturnRows(rows)

	result, total, err := repo.List(ctx, offset, limit)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
	assert.Equal(t, users[0].ID, result[0].ID)
	assert.Equal(t, users[1].ID, result[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepositoryExistsByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupUserRepositoryTest(t)
	defer cleanup()

	repo := persistence.NewUserRepository(gormDB)
	ctx := context.Background()

	email := "test@example.com"

	t.Run("exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE email = \\? AND `users`.`deleted_at` IS NULL").
			WithArgs(email).
			WillReturnRows(rows)

		exists, err := repo.ExistsByEmail(ctx, email)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM `users` WHERE email = \\? AND `users`.`deleted_at` IS NULL").
			WithArgs(email).
			WillReturnRows(rows)

		exists, err := repo.ExistsByEmail(ctx, email)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
