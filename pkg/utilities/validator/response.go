package validator

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo"
)

type ResponseMessage struct {
	Code        int         `json:"code,omitempty"`
	MessageCode string      `json:"messageCode,omitempty"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}
type Response struct {
	echo.Context
	Message func(string) ResponseMessage
}

func (e Response) HandleError(err error) error {
	errMessage := e.Message(err.Error())
	if reflect.ValueOf(errMessage).IsZero() {
		errMessage := struct {
			Code    int    `json:"code,omitempty"`
			Message string `json:"message,omitempty"`
		}{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}

		return e.JSON(http.StatusInternalServerError, errMessage)
	}

	return e.JSON(http.StatusInternalServerError, errMessage)
}

func (e Response) HandleSuccess(response interface{}) error {
	return e.JSON(http.StatusOK, response)
}

func (e Response) HandleCreated(response interface{}) error {
	return e.JSON(http.StatusCreated, response)
}

func (e Response) HandleBadRequest(err error) error {
	return e.JSON(http.StatusBadRequest, FormatValidationErrors(err))
}

func (e Response) NotFound(err error) error {
	message := e.Message(err.Error())
	if reflect.ValueOf(message).IsZero() {
		return e.JSON(http.StatusNotFound, err)
	}
	return e.JSON(http.StatusNotFound, message)
}
