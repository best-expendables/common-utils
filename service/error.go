package service

import (
	"fmt"
	"strings"

	"github.com/best-expendables/common-utils/util/validation"
)

type ValidationError map[string][]string

func NewValidationError(err error) ValidationError {
	return validation.ParseValidationErr(err)
}

func (m ValidationError) Error() string {
	result := make([]string, 0)
	for key, values := range m {
		result = append(result, fmt.Sprintf("%v: %v", key, strings.Join(values, ", ")))
	}
	return strings.Join(result, ", ")
}

type ForbiddenError string

func (f ForbiddenError) Error() string {
	return string(f)
}

type Unauthorized string

func (f Unauthorized) Error() string {
	return string(f)
}

type NotFoundError string

func (f NotFoundError) Error() string {
	return string(f)
}

type InternalServerError struct {
	Err error
}

func (i InternalServerError) Error() string {
	return i.Err.Error()
}

type BadRequestError string

func (f BadRequestError) Error() string {
	return string(f)
}
