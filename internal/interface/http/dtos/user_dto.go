package dtos

// CreateUserRequest represents payload to create user.
type CreateUserRequest struct {
	Name        string `json:"name" binding:"required"`
	DeviceID    string `json:"device_id" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	BackupPhone string `json:"backup_phone"`
}

// UserResponse represents user output.
type UserResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DeviceID    string `json:"device_id"`
	Phone       string `json:"phone"`
	BackupPhone string `json:"backup_phone"`
}
