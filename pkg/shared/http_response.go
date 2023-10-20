package shared

// Response standard response for handle http request
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// SuccessResponse to handle success response
func SuccessResponse(data any) Response {
	return Response{
		Message: "success",
		Data:    data,
	}
}

// ErrorResponse to handle error response
func ErrorResponse(err string) Response {
	return Response{
		Message: err,
	}
}
