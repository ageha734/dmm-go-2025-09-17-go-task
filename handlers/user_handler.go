package handlers

import (
	"net/http"

	"github.com/dmm-com/dmm-go-2025-09-17-go-task/database"
	"github.com/dmm-com/dmm-go-2025-09-17-go-task/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ErrUserNotFound = "User not found"
)

type DBProvider interface {
	GetDB() *gorm.DB
}

type DefaultDBProvider struct{}

func (d *DefaultDBProvider) GetDB() *gorm.DB {
	return database.GetDB()
}

var dbProvider DBProvider = &DefaultDBProvider{}

func SetDBProvider(provider DBProvider) {
	dbProvider = provider
}

func GetUsers(c *gin.Context) {
	var users []models.User
	db := dbProvider.GetDB()

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db := dbProvider.GetDB()

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty"})
		return
	}

	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	db := dbProvider.GetDB()
	if err := db.Create(&user).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" ||
			err.Error() == "Error 1062: Duplicate entry" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	db := dbProvider.GetDB()

	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Age != 0 {
		user.Age = req.Age
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db := dbProvider.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound})
		return
	}

	if err := db.Delete(&user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func GetUserStats(c *gin.Context) {
	db := dbProvider.GetDB()
	var totalUsers int64
	var avgAge *float64

	if err := db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}

	if err := db.Model(&models.User{}).Select("AVG(age)").Where("age IS NOT NULL").Scan(&avgAge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate average age"})
		return
	}

	var averageAge float64
	if avgAge != nil {
		averageAge = *avgAge
	} else {
		averageAge = 0
	}

	stats := gin.H{
		"total_users": totalUsers,
		"average_age": averageAge,
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API is running successfully",
	})
}
