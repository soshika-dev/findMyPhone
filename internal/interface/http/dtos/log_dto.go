package dtos

// CreateLogRequest payload.
type CreateLogRequest struct {
	DeviceID  string  `json:"device_id" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

// LogResponse output.
type LogResponse struct {
	ID        uint    `json:"id"`
	DeviceID  string  `json:"device_id"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	CreatedAt string  `json:"created_at"`
}
