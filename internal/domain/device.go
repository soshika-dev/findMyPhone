package domain

import "time"

// Device represents a trackable device entry.
type Device struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DeviceID   string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"device_id"`
	IMEI       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"imei"`
	Generation string    `gorm:"type:varchar(255)" json:"generation"`
	Name       string    `gorm:"type:varchar(255)" json:"name"`
	Lost       bool      `gorm:"default:false" json:"lost"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
