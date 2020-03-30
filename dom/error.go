package dom

// Error represents ML Proxy Error
type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

// NewError creates ML Proxy Error
func NewError(err string, message string, status int) *Error {
	return &Error{
		Message: message,
		Status:  status,
		Error:   err,
	}
}

// UnknownError creates ML Proxy Error for unknown cases
func UnknownError() *Error {
	return &Error{
		Status:  500,
		Message: "Unknown",
		Error:   "Unknown",
	}
}
