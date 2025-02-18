package models

import "time"

type User struct {
	ID        uint    `gorm:"primaryKey"`
	Username  *string `json:"username,omitempty"` // Теперь это указатель, может быть nil
	Email     string  `gorm:"unique;not null"`
	Password  string  `gorm:"not null"`
	Role      string  `gorm:"not null;default:user"`
	CreatedAt time.Time
}
