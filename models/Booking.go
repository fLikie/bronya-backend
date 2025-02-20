package models

import "time"

type Booking struct {
	ID        uint      `json:"id"`
	PlaceID   uint      `json:"place_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	TimeSlot  string    `json:"time_slot" gorm:"not null"` // Формат HH:MM
	CreatedAt time.Time `json:"created_at"`
}
