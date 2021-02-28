package response

import (
	"net/http"

	"github.com/best-expendables/common-utils/service"
)

func ConvertServiceError(err error) ApiResponse {
	switch err.(type) {
	case service.ValidationError:
		return CreateValidationErrResponse(err.(service.ValidationError))
	case service.ForbiddenError:
		return ErrorResponse(err, http.StatusForbidden)
	case service.NotFoundError:
		return ErrorResponse(err, http.StatusNotFound)
	case service.Unauthorized:
		return ErrorResponse(err, http.StatusUnauthorized)
	case service.BadRequestError:
		return ErrorResponse(err, http.StatusBadRequest)
	}
	return ErrorResponse(err, http.StatusInternalServerError)
}
