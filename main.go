package main

import (
	"bronya/database"
	"bronya/middlewares"
	"bronya/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	err = connectDatabase.AutoMigrate(&models.User{}, &models.Place{}, &models.Booking{})
	if err != nil {
		return
	}
	db = connectDatabase
}

func CreatePlace(c *gin.Context) {
	name := c.PostForm("name")
	location := c.PostForm("location")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload required"})
		return
	}

	// Сохраняем файл в папку uploads
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "/var/www/bronya-web/uploads/" + filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Создаём запись в базе
	place := models.Place{Name: name, Location: location, Image: filename}
	if err := db.Create(&place).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create place"})
		return
	}

	c.JSON(http.StatusOK, place)
}

func CreateBooking(c *gin.Context) {
	var input struct {
		PlaceID  uint   `json:"place_id" binding:"required"`
		TimeSlot string `json:"time_slot" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	booking := models.Booking{
		PlaceID:  input.PlaceID,
		UserID:   userID.(uint),
		TimeSlot: input.TimeSlot,
	}

	if err := db.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func GetBookingsForPlace(c *gin.Context) {
	placeID := c.Param("id")

	var bookings []models.Booking
	if err := db.Where("place_id = ?", placeID).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func makeAdmin(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
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
	r.POST("/api/register", middlewares.RegisterHandler)
	r.POST("/api/login", middlewares.LoginHandler)
	r.GET("/api/profile", middlewares.AuthChecking(), Profile)
	r.GET("/api/places", middlewares.AuthChecking(), GetPlaces)
	r.POST("/api/places", middlewares.AuthChecking(), middlewares.AdminChecking(), CreatePlace)
	r.POST("/api/bookings", middlewares.AuthChecking(), CreateBooking)
	r.POST("/api/make-admin", middlewares.AuthChecking(), middlewares.AdminChecking(), makeAdmin)
	r.GET("/api/users", middlewares.AuthChecking(), middlewares.AdminChecking(), GetUsers)
	r.GET("/api/places/:id", middlewares.AuthChecking(), GetPlace)
	r.PUT("/api/places/:id", middlewares.AuthChecking(), UpdatePlace)
	r.GET("/api/places/:id/bookings", middlewares.AuthChecking(), GetBookingsForPlace)
	r.Run(":8080")
}

func Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var user models.User
	db.First(&user, userID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetPlaces(c *gin.Context) {
	var places []models.Place
	if err := db.Find(&places).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch places"})
		return
	}
	c.JSON(http.StatusOK, places)
}

func GetPlace(c *gin.Context) {
	var place models.Place
	placeID := c.Param("id")

	if err := db.First(&place, placeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}
	c.JSON(http.StatusOK, place)
}

func UpdatePlace(c *gin.Context) {
	var place models.Place
	placeID := c.Param("id")

	if err := db.First(&place, placeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}

	name := c.PostForm("name")
	location := c.PostForm("location")

	// Обновляем текстовые данные
	place.Name = name
	place.Location = location

	// Проверяем, есть ли новое изображение
	file, err := c.FormFile("image")
	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filepath := "/var/www/bronya-web/uploads/" + filename

		// Сохраняем новый файл
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}

		// Удаляем старое изображение
		oldFilePath := "/var/www/bronya-web/uploads/" + place.Image
		os.Remove(oldFilePath)

		// Сохраняем новый путь
		place.Image = filename
	}

	if err := db.Save(&place).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update place"})
		return
	}

	c.JSON(http.StatusOK, place)
}
