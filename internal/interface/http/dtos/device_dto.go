package dtos

// CreateDeviceRequest represents payload to register device.
type CreateDeviceRequest struct {
	DeviceID   string `json:"device_id" binding:"required"`
	IMEI       string `json:"imei" binding:"required"`
	Generation string `json:"generation"`
	Name       string `json:"name"`
	Lost       bool   `json:"lost"`
}

// DeviceResponse is output for device.
type DeviceResponse struct {
	ID         uint   `json:"id"`
	DeviceID   string `json:"device_id"`
	IMEI       string `json:"imei"`
	Generation string `json:"generation"`
	Name       string `json:"name"`
	Lost       bool   `json:"lost"`
}
