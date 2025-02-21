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

func GetUserById(user *models.User, userId uint) error {
	return DB.First(&user, userId).Error
}

func GetAllUsers(users *[]models.User) error {
	return DB.Find(users).Error
}

func CreateBooking(booking *models.Booking) error {
	return DB.Create(&booking).Error
}

func GetBookingsForPlace(bookings *[]models.Booking, placeId string) error {
	return DB.Where("place = ?", placeId).First(&bookings).Error
}

func CreatePlace(place *models.Place) error {
	return DB.Create(&place).Error
}

func GetPlace(place *models.Place, placeId string) error {
	return DB.First(&place, placeId).Error
}

func GetPlaces(places *[]models.Place) error {
	return DB.Find(&places).Error
}

func UpdatePlace(place *models.Place) error {
	return DB.Save(&place).Error
}

func UpdateUser(user *models.User) error {
	return DB.Save(&user).Error
}
