package dom

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func UnknownError() *Error {
	return &Error{
		Status:  500,
		Message: "Unknown",
		Error:   "Unknown",
	}
}
