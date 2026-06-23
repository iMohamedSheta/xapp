package enums

import (
	"net/http"

	"github.com/iMohamedSheta/xerr"
)

/*
| This file is to set custom errors types for the application
*/
const (
	XErrValidationError xerr.ErrorType = iota + 1000
	XErrBadRequestError
	XErrBadRequestBindingError
	XErrUnAuthorizedError
	XErrForbiddenError
	XErrNotFoundError
	XErrServerError
)

type ErrorCode string

const (
	ErrCodeInternalError   ErrorCode = "internal_error"
	ErrCodeValidationError ErrorCode = "validation_error"
	ErrCodeUnauthorized    ErrorCode = "unauthorized_error"
	ErrCodeForbidden       ErrorCode = "forbidden_error"
	ErrCodeNotFound        ErrorCode = "not_found_error"
	ErrCodeBadRequest      ErrorCode = "bad_request_error"
)

func (e ErrorCode) String() string {
	return string(e)
}

// Map ErrorCode → HTTP Status
func (e ErrorCode) StatusCode() int {
	return map[ErrorCode]int{
		ErrCodeInternalError:   http.StatusInternalServerError,
		ErrCodeValidationError: http.StatusUnprocessableEntity,
		ErrCodeUnauthorized:    http.StatusUnauthorized,
		ErrCodeForbidden:       http.StatusForbidden,
		ErrCodeNotFound:        http.StatusNotFound,
		ErrCodeBadRequest:      http.StatusBadRequest,
	}[e]
}

// ErrorCode default messages (public)
func (e ErrorCode) Message() string {
	return map[ErrorCode]string{
		ErrCodeInternalError:   "Something went wrong, please try again later.",
		ErrCodeValidationError: "Some fields are invalid. Please check your input.",
		ErrCodeUnauthorized:    "You must be logged in to access this resource.",
		ErrCodeForbidden:       "You don’t have permission to perform this action.",
		ErrCodeNotFound:        "The requested resource could not be found.",
		ErrCodeBadRequest:      "The request could not be understood by the server.",
	}[e]
}

func (e ErrorCode) MessageAr() string {
	return map[ErrorCode]string{
		ErrCodeInternalError:   "حدث خطأ ما، يرجى المحاولة لاحقًا.",
		ErrCodeValidationError: "بعض الحقول غير صالحة. يرجى مراجعة الإدخال.",
		ErrCodeUnauthorized:    "يجب تسجيل الدخول للوصول إلى هذا المورد.",
		ErrCodeForbidden:       "ليس لديك صلاحية لتنفيذ هذا الإجراء.",
		ErrCodeNotFound:        "لم يتم العثور على المورد المطلوب.",
		ErrCodeBadRequest:      "لم يتمكن الخادم من فهم هذا الطلب.",
	}[e]
}

func GetErrorCode(t xerr.ErrorType) ErrorCode {
	switch t {
	case XErrValidationError:
		return ErrCodeValidationError
	case XErrBadRequestError, XErrBadRequestBindingError:
		return ErrCodeBadRequest
	case XErrUnAuthorizedError:
		return ErrCodeUnauthorized
	case XErrForbiddenError:
		return ErrCodeForbidden
	case XErrNotFoundError:
		return ErrCodeNotFound
	default:
		return ErrCodeInternalError
	}
}
