package main

import (
	"bronya/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	gorm "gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectDatabase := database.ConnectDatabase()
	err = connectDatabase.AutoMigrate(&User{}, &Place{}, &Booking{})
	if err != nil {
		return
	}
	db = connectDatabase
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func createPlace(c *gin.Context) {
	var place Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&place)
	c.JSON(http.StatusOK, gin.H{"message": "Place added successfully"})
}

func createBooking(c *gin.Context) {
	var booking Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&booking)
	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
}

func makeAdmin(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = "admin"
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}

func main() {
	r := gin.Default()
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)
	r.GET("/profile", AuthMiddleware(), Profile)
	r.POST("/places", AuthMiddleware(), adminMiddleware(), createPlace)
	r.POST("/bookings", AuthMiddleware(), createBooking)
	r.POST("/make-admin", AuthMiddleware(), adminMiddleware(), makeAdmin)
	r.Run(":8080")
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null;default:user"`
	CreatedAt time.Time
}

type Place struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Location  string `gorm:"not null"`
	CreatedAt time.Time
}

type Booking struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	PlaceID   uint      `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
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
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	db.Create(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
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

func CreatePlace(c *gin.Context) {
	var place Place
	if err := c.ShouldBindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&place)
	c.JSON(http.StatusOK, gin.H{"message": "Place added successfully"})
}

func CreateBooking(c *gin.Context) {
	var booking Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&booking)
	c.JSON(http.StatusOK, gin.H{"message": "Booking created successfully"})
}
