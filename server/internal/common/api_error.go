package common

import (
	"net/http"
	"strings"
)

type ApiError struct {
	Code          int               `json:"code"`
	ErrorsDetails []ApiErrorDetails `json:"errors"`
}

func (apiError *ApiError) Error() string {
	var messageBuilder strings.Builder

	for idx, apiErrorDetails := range apiError.ErrorsDetails {
		messageBuilder.WriteString(apiErrorDetails.Error())

		isLastError := (idx == len(apiError.ErrorsDetails)-1)
		if !isLastError {
			messageBuilder.WriteString(" | ")
		}
	}
	return messageBuilder.String()
}

func NewApiErrors(code int, errors []ApiErrorDetails) *ApiError {
	return &ApiError{
		Code:          code,
		ErrorsDetails: errors,
	}
}

func NewBadRequestError(apiErrorDetails ApiErrorDetails) *ApiError {
	return NewApiErrors(http.StatusBadRequest, []ApiErrorDetails{apiErrorDetails})
}

func NewConflictError(apiErrorDetails ApiErrorDetails) *ApiError {
	return NewApiErrors(http.StatusBadRequest, []ApiErrorDetails{apiErrorDetails})
}

func NewUnprocessableEntityError(apiErrorDetails ApiErrorDetails) *ApiError {
	return NewApiErrors(http.StatusUnprocessableEntity, []ApiErrorDetails{apiErrorDetails})
}

func NewUnprocessableEntityErrors(apiErrorsDetails []ApiErrorDetails) *ApiError {
	return NewApiErrors(http.StatusUnprocessableEntity, apiErrorsDetails)
}

func NewUnauthorizedError(message string) *ApiError {
	apiError := ApiErrorDetails{
		Message: message,
	}
	return NewApiErrors(http.StatusUnauthorized, []ApiErrorDetails{apiError})
}

func NewInternalServerError(message string) *ApiError {
	apiError := ApiErrorDetails{
		Message: message,
	}
	return NewApiErrors(http.StatusInternalServerError, []ApiErrorDetails{apiError})
}
