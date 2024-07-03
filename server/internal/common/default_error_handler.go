package common

import (
	"net/http"

	"github.com/labstack/echo"
)

type httpErrorHandler struct{}

func NewHttpErrorHandler() *httpErrorHandler {
	return &httpErrorHandler{}
}

func (errorHandler *httpErrorHandler) Handler(err error, c echo.Context) {
	var handledApiError *ApiError

	if apiError, ok := err.(*ApiError); ok {

		handledApiError = apiError
	} else {
		code := http.StatusInternalServerError
		message := err.Error()

		httpErr, ok := err.(*echo.HTTPError)
		if ok {
			code = httpErr.Code
			message = httpErr.Message.(string)
		}

		handledApiError = &ApiError{
			Code: code,
			ErrorsDetails: []ApiErrorDetails{
				{
					Message: message,
				},
			},
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(handledApiError.Code)
		} else {
			err = c.JSON(handledApiError.Code, map[string][]ApiErrorDetails{
				"errors": handledApiError.ErrorsDetails,
			})
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}
