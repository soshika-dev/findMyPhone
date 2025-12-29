package dtos

// Response defines standard API response envelope.
type Response struct {
	Data    any    `json:"data,omitempty"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
