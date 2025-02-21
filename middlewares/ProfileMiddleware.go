package middlewares

import (
	"bronya/database"
	"bronya/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var user models.User
	err := database.GetUserById(&user, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.GetAllUsers(&users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func MakeAdmin(c *gin.Context) {
	var input struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.FindUser(&user, input.Phone); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = "admin"
	err := database.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}
