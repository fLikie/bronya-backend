package models

import "time"

type Booking struct {
	ID        uint      `json:"id"`
	PlaceID   uint      `json:"place_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Date      string    `json:"date" gorm:"not null"` // Дата в формате YYYY-MM-DD
	TimeSlot  string    `json:"time_slot" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
