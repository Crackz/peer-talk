package common

import (
	"fmt"
)

type ApiErrorDetails struct {
	Message string `json:"message"`
	Param   string `json:"param,omitempty"`
}

func (errorDetails *ApiErrorDetails) Error() string {
	return fmt.Sprintf("message: %v, param: %v", errorDetails.Message, errorDetails.Param)
}
