package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ageha734/dmm-go-2025-09-17-go-task/handlers"
	"github.com/ageha734/dmm-go-2025-09-17-go-task/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	usersEndpoint    = "/users"
	userByIDEndpoint = "/users/:id"
	testUserName     = "テストユーザー"
	testUserEmail    = "test@example.com"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("failed to migrate database")
	}

	return db
}

type TestDBProvider struct {
	db *gorm.DB
}

func (t *TestDBProvider) GetDB() *gorm.DB {
	return t.db
}

func TestHealthCheck(t *testing.T) {
	router := gin.New()
	router.GET("/health", handlers.HealthCheck)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "API is running successfully", response["message"])
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := setupTestDB()
	handlers.SetDBProvider(&TestDBProvider{db: testDB})

	router.POST(usersEndpoint, handlers.CreateUser)

	user := models.CreateUserRequest{
		Name:  testUserName,
		Email: testUserEmail,
		Age:   25,
	}

	jsonData, _ := json.Marshal(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a map[string]interface{}")
	assert.Equal(t, testUserName, data["name"])
	assert.Equal(t, testUserEmail, data["email"])
	assert.Equal(t, float64(25), data["age"])
}

func TestCreateUserInvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := setupTestDB()
	handlers.SetDBProvider(&TestDBProvider{db: testDB})

	router.POST(usersEndpoint, handlers.CreateUser)

	invalidUser := map[string]interface{}{
		"email": "invalid-email",
		"age":   -1,
	}

	jsonData, _ := json.Marshal(invalidUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", usersEndpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "required")
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := setupTestDB()

	testUser := models.User{
		Name:  testUserName,
		Email: testUserEmail,
		Age:   25,
	}
	testDB.Create(&testUser)

	handlers.SetDBProvider(&TestDBProvider{db: testDB})

	router.GET(usersEndpoint, handlers.GetUsers)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", usersEndpoint, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].([]interface{})
	assert.True(t, ok, "data should be a []interface{}")
	assert.Len(t, data, 1)

	user, ok := data[0].(map[string]interface{})
	assert.True(t, ok, "user should be a map[string]interface{}")
	assert.Equal(t, testUserName, user["name"])
	assert.Equal(t, testUserEmail, user["email"])
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := setupTestDB()

	testUser := models.User{
		Name:  testUserName,
		Email: testUserEmail,
		Age:   25,
	}
	testDB.Create(&testUser)

	handlers.SetDBProvider(&TestDBProvider{db: testDB})

	router.GET(userByIDEndpoint, handlers.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a map[string]interface{}")
	assert.Equal(t, testUserName, data["name"])
	assert.Equal(t, testUserEmail, data["email"])
}

func TestGetUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	testDB := setupTestDB()
	handlers.SetDBProvider(&TestDBProvider{db: testDB})

	router.GET(userByIDEndpoint, handlers.GetUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}
