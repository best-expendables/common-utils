package response

import (
	"errors"
	"net/http"

	"github.com/best-expendables/common-utils/util/validation"
	"gorm.io/gorm"
)

type ApiResponse struct {
	Code    int               `json:"-"`
	Data    interface{}       `json:"data,omitempty"`
	Headers map[string]string `json:"-"`
	Errors  interface{}       `json:"errors,omitempty"`
}

type Error struct {
	Message    string              `json:"message,omitempty"`
	Code       int                 `json:"code,omitempty"`
	StatusCode int                 `json:"status_code,omitempty"`
	Errors     map[string][]string `json:"errors,omitempty"`
}

// Ok 200-OK with data in body
func Ok(data interface{}) ApiResponse {
	return ApiResponse{http.StatusOK, data, nil, nil}
}

// Accepted 202-Accepted with data in body
func Accepted(data interface{}) ApiResponse {
	return ApiResponse{http.StatusAccepted, data, nil, nil}
}

// Created 201-Created with data in body
func Created(data interface{}) ApiResponse {
	return ApiResponse{http.StatusCreated, data, nil, nil}
}

// PartialContent 206-Partial Content with data in body
func PartialContent(data interface{}) ApiResponse {
	return ApiResponse{http.StatusPartialContent, data, nil, nil}
}

// BadRequest 400-Bad Request
func BadRequest(err error) ApiResponse {
	return ErrorResponse(err, http.StatusBadRequest)
}

// 500 Internal server error
func InternalServerError(err error) ApiResponse {
	return ErrorResponse(err, http.StatusInternalServerError)
}

// ValidationError 422-Unprocessable Enttity
func ValidationError(err error) ApiResponse {
	errMessages := validation.ParseValidationErr(err)
	var response ApiResponse
	response.Code = http.StatusUnprocessableEntity
	response.Errors = Error{
		Message:    http.StatusText(http.StatusUnprocessableEntity),
		StatusCode: http.StatusUnprocessableEntity,
		Code:       http.StatusUnprocessableEntity,
		Errors:     errMessages,
	}
	return response
}

// ErrorResponse Custom error response
func ErrorResponse(err error, status int) ApiResponse {
	var response ApiResponse
	response.Code = status
	response.Errors = Error{
		Message:    err.Error(),
		StatusCode: status,
		Code:       status,
	}
	return response
}

func RecordNotFoundError(err error, message string) ApiResponse {
	if err == gorm.ErrRecordNotFound {
		return ErrorResponse(errors.New(message), http.StatusNotFound)
	}

	return ErrorResponse(err, http.StatusInternalServerError)
}

func CreateValidationErrResponse(errMessages map[string][]string) ApiResponse {
	var response ApiResponse
	response.Code = http.StatusUnprocessableEntity
	response.Errors = Error{
		Message:    http.StatusText(http.StatusUnprocessableEntity),
		StatusCode: http.StatusUnprocessableEntity,
		Code:       http.StatusUnprocessableEntity,
		Errors:     errMessages,
	}
	return response
}
