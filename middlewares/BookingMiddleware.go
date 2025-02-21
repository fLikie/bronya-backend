package middlewares

import (
	"bronya/database"
	"bronya/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateBooking(c *gin.Context) {
	var input struct {
		PlaceID  uint   `json:"place_id" binding:"required"`
		Date     string `json:"date" binding:"required"`
		TimeSlot string `json:"time_slot" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Invalid input error:", err.Error()) // Логируем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Printf("Booking request: PlaceID=%d, Date=%s, TimeSlot=%s, UserID=%d", input.PlaceID, input.Date, input.TimeSlot, userID.(uint))

	booking := models.Booking{
		PlaceID:  input.PlaceID,
		UserID:   userID.(uint),
		Date:     input.Date,
		TimeSlot: input.TimeSlot,
	}

	if err := database.CreateBooking(&booking); err != nil {
		log.Println("DB error:", err.Error()) // Логируем ошибку базы данных
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func GetBookingsForPlace(c *gin.Context) {
	placeID := c.Param("id")

	var bookings []models.Booking
	if err := database.GetBookingsForPlace(&bookings, placeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}
