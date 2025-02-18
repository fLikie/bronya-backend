package models

import "time"

type Place struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Location  string `gorm:"not null"`
	CreatedAt time.Time
}
