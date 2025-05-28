package validator

import (
	"net/http"
)

const (
	ERR_VALIDATION_FAILED    = `VALIDATION_FAILED`
	ERR_INVALID_EMAIL_FORMAT = `INVALID_EMAIL_FORMAT`
)

var ErrorMessage = map[string]string{
	`VALIDATION_FAILED`:    `Failed for validate the input`,
	`INVALID_EMAIL_FORMAT`: `Invalid email format provided`,
}

var ErrorStatusCode = map[string]int{
	`VALIDATION_FAILED`:    http.StatusNotAcceptable,
	`INVALID_EMAIL_FORMAT`: http.StatusBadRequest,
}
