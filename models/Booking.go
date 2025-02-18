package models

import "time"

type Booking struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	PlaceID   uint      `gorm:"not null"`
	Date      time.Time `gorm:"not null"`
	CreatedAt time.Time
}
