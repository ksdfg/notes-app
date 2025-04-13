package utils

// ApiResponse is a struct that defines the schema for all API responses.
type ApiResponse struct {
	// Success is a boolean indicating whether the API call was successful or not.
	Success bool `json:"success"`

	// Message is a human-readable message indicating the result of the API call. If the call was unsuccessful, this
	//message should contain an error message.
	Message string `json:"message"`
}
