package main

import (
	"bronya/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	gorm "gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

var db *gorm.DB

func main() {
	// Подключаем базу
	db = database.ConnectDatabase()

	// Создаём API сервер
	r := gin.Default()

	// Тестовый маршрут
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Сервис бронирования работает! 🚀"})
	})

	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/profile", AuthMiddleware(), Profile)

	// Запускаем сервер
	r.Run(":8080")
}

type User struct {
	ID        int    `gorm:"primary_key"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
}

func registerHandler(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func loginHandler(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var user User

	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var user User
	db.First(&user, userID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Next()
	}
}
