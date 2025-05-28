package responses

import "net/http"

type Response struct {
	Code    string      `json:"code" example:"SUCCESS"`
	Message string      `json:"message" example:"successfully"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code  string      `json:"code" example:"UNHANDLED_EXCEPTION"`
	Error ErrorDetail `json:"error"`
}
type ErrorDetail struct {
	Message string `json:"message" example:"Internal Server Error"`
	Stack   string `json:"stack,omitempty" example:"Error:Database error"`
}

func Ok(okCode int, message string, payload interface{}) *Response {
	return &Response{
		Code:    StatusBusinessCode(okCode),
		Message: message,
		Data:    payload,
	}
}

func Error(errorCode int, stack string) *ErrorResponse {
	return &ErrorResponse{
		Code: StatusBusinessCode(errorCode),
		Error: ErrorDetail{
			Message: http.StatusText(errorCode),
			Stack:   "Error: " + stack,
		},
	}
}
