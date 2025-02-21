package middlewares

import (
	"bronya/database"
	"bronya/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

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
	if err := database.CreatePlace(&place); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create place"})
		return
	}

	c.JSON(http.StatusOK, place)
}

func GetPlaces(c *gin.Context) {
	var places []models.Place
	if err := database.GetPlaces(&places); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch places"})
		return
	}
	c.JSON(http.StatusOK, places)
}

func GetPlace(c *gin.Context) {
	var place models.Place
	placeID := c.Param("id")

	if err := database.GetPlace(&place, placeID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Place not found"})
		return
	}
	c.JSON(http.StatusOK, place)
}

func UpdatePlace(c *gin.Context) {
	var place models.Place
	placeID := c.Param("id")

	if err := database.GetPlace(&place, placeID); err != nil {
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
		err := os.Remove(oldFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image"})
			return
		}

		// Сохраняем новый путь
		place.Image = filename
	}

	if err := database.UpdatePlace(&place); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update place"})
		return
	}

	c.JSON(http.StatusOK, place)
}
