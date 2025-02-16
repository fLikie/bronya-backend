package main

import (
	"bronya/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Подключаем базу
	database.ConnectDatabase()

	// Создаём API сервер
	r := gin.Default()

	// Тестовый маршрут
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Сервис бронирования работает! 🚀"})
	})

	// Запускаем сервер
	r.Run(":8080")
}
