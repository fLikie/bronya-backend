package database

import (
	"bronya/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// Глобальная переменная для базы
var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используем переменные среды")
	}

	// Формируем строку подключения
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Подключаемся к базе
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	fmt.Println("✅ Успешное подключение к базе данных")

	return DB
}

func CreateUser(user *models.User) {
	DB.Create(&user)
}

func FindUser(user *models.User, phone string) error {
	return DB.Where("phone = ?", phone).First(&user).Error
}
