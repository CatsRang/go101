package example

// CreateRequest represents the payload for creating a resource
type CreateRequest struct {
	Name     string `json:"name"`
	Duration int    `json:"duration_ms"` // Simulate work duration
}

// CreateResponse represents the response
type CreateResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
