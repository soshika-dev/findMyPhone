package domain

import "time"

// Log captures location ping of a device.
type Log struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DeviceID  string    `gorm:"type:varchar(255);not null;index" json:"device_id"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
