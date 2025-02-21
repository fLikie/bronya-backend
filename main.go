package main

import (
	"bronya/database"
	"bronya/middlewares"
	"bronya/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	_ "net/http"
)

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
}

func main() {
	r := gin.Default()
	// auth api
	r.POST("/api/register", middlewares.RegisterHandler)
	r.POST("/api/login", middlewares.LoginHandler)

	// profile api
	r.GET("/api/profile", middlewares.AuthChecking(), middlewares.Profile)

	// places api
	r.GET("/api/places", middlewares.AuthChecking(), middlewares.GetPlaces)
	r.POST("/api/places", middlewares.AuthChecking(), middlewares.AdminChecking(), middlewares.CreatePlace)
	r.GET("/api/places/:id", middlewares.AuthChecking(), middlewares.GetPlace)
	r.PUT("/api/places/:id", middlewares.AuthChecking(), middlewares.UpdatePlace)

	// bookings api
	r.POST("/api/bookings", middlewares.AuthChecking(), middlewares.CreateBooking)
	r.GET("/api/places/:id/bookings", middlewares.AuthChecking(), middlewares.GetBookingsForPlace)

	// users control api
	r.POST("/api/make-admin", middlewares.AuthChecking(), middlewares.AdminChecking(), middlewares.MakeAdmin)
	r.GET("/api/users", middlewares.AuthChecking(), middlewares.AdminChecking(), middlewares.GetUsers)

	// run listening
	r.Run(":8080")
}
