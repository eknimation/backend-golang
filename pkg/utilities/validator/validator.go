package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Email validation regex pattern - more comprehensive than default
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

type fieldError struct {
	err validator.FieldError
}

type ValidationError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

/*
 *	Validate
 */
func Validate(input interface{}) error {
	validate = validator.New()

	_ = validate.RegisterValidation("isComplexPassword", validateComplexPassword)
	_ = validate.RegisterValidation("emailFormat", validateEmailFormat)
	return validate.Struct(input)
}

/*
 *	Format validation rrrors
 */
func FormatValidationErrors(err error) ErrResponse {
	errs := []ValidationError{}
	for _, err := range err.(validator.ValidationErrors) {
		validationErr := ValidationError{
			Code:    strings.ToUpper(err.ActualTag()),
			Field:   err.Field(),
			Message: fieldError{err}.String(),
		}

		errs = append(errs, validationErr)
	}

	serializedErr := ErrResponse{
		Code:    ErrorStatusCode[ERR_VALIDATION_FAILED],
		Message: ErrorMessage[ERR_VALIDATION_FAILED],
		Errors:  errs,
	}

	return *(*ErrResponse)(unsafe.Pointer(&serializedErr))
}

func (q fieldError) String() string {
	var sb strings.Builder

	sb.WriteString("Validation failed on field '" + q.err.Field() + "'")
	sb.WriteString(", condition: " + q.err.ActualTag())

	if q.err.Param() != "" {
		sb.WriteString(" { " + q.err.Param() + " }")
	}

	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", q.err.Value()))
	}

	return sb.String()
}

// function convert string to time.Time
func ConvertStringToTime(date string) (time.Time, error) {
	var response error
	result, err := time.Parse(time.RFC3339, date)
	if err != nil {
		response = errors.New("Date is wrong format")
		return result, response
	}

	return result, response
}

// validate date range by start date and end date
func ValidateDateRange(startDate time.Time, endDate time.Time) error {
	var response error

	// end date can equal to start date but not before
	if endDate.Before(startDate) {
		response = errors.New("End date must be after start date")
		return response
	}

	return response
}

func validateComplexPassword(fl validator.FieldLevel) bool {
	var (
		hasMinLen  = false
		hasMaxLen  = true
		hasNumber  = false
		hasUpper   = false
		hasLower   = false
		hasSpecial = false
	)
	pass := fl.Field().String()

	if len(pass) >= 8 {
		hasMinLen = true
	}

	if len(pass) > 32 {
		hasMaxLen = false
	}

	for _, char := range pass {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case strings.ContainsRune("!@#$%^&*()-_+=", char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasMaxLen && hasNumber && hasUpper && hasLower && hasSpecial
}

// validateEmailFormat provides enhanced email validation
func validateEmailFormat(fl validator.FieldLevel) bool {
	email := strings.TrimSpace(fl.Field().String())

	// Check if email is empty
	if email == "" {
		return false
	}

	// Check email length (RFC 5321 limits)
	if len(email) > 254 {
		return false
	}

	// Use regex pattern for validation
	if !emailRegex.MatchString(email) {
		return false
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Local part checks
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Domain part checks
	if len(domainPart) == 0 || len(domainPart) > 253 {
		return false
	}

	// Domain must have at least one dot
	if !strings.Contains(domainPart, ".") {
		return false
	}

	// Domain cannot start or end with a dot or hyphen
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") ||
		strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return false
	}

	return true
}

// IsValidEmail validates email format independently without using struct validation
// This can be used for standalone email validation
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	// Check if email is empty
	if email == "" {
		return false
	}

	// Check email length (RFC 5321 limits)
	if len(email) > 254 {
		return false
	}

	// Use regex pattern for validation
	if !emailRegex.MatchString(email) {
		return false
	}

	// Additional checks
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Local part checks
	if len(localPart) == 0 || len(localPart) > 64 {
		return false
	}

	// Domain part checks
	if len(domainPart) == 0 || len(domainPart) > 253 {
		return false
	}

	// Domain must have at least one dot
	if !strings.Contains(domainPart, ".") {
		return false
	}

	// Domain cannot start or end with a dot or hyphen
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") ||
		strings.HasPrefix(domainPart, "-") || strings.HasSuffix(domainPart, "-") {
		return false
	}

	return true
}
