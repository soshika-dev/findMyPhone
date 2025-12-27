package domain

import "time"

// User represents an account tied to a device.
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	DeviceID    string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"device_id"`
	Phone       string    `gorm:"type:varchar(50);not null" json:"phone"`
	BackupPhone string    `gorm:"type:varchar(50)" json:"backup_phone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
