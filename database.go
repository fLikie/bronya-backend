package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Глобальная переменная для базы
var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем переменные среды")
	}

	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"45.143.95.91",
		"booking_user",
		"password",
		"booking_service",
		"5432",
	)

	// Подключаемся к базе
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	fmt.Println("✅ Успешное подключение к базе данных")
}
