package common

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	notStandardValidators "github.com/go-playground/validator/v10/non-standard/validators"

	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type DefaultValidator struct {
	validator    *validator.Validate
	translations ut.Translator
}

func NewDefaultValidator() *DefaultValidator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	newValidator := validator.New(validator.WithRequiredStructEnabled())
	newValidator.RegisterValidation("notblank", notStandardValidators.NotBlank)
	newValidator.RegisterTranslation("notblank", trans, func(ut ut.Translator) error {
		return ut.Add("notblank", "{0} can't be blank value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notblank", fe.Field())

		return t
	})
	newValidator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	en_translations.RegisterDefaultTranslations(newValidator, trans)

	return &DefaultValidator{
		validator:    newValidator,
		translations: trans,
	}
}

func (dv *DefaultValidator) Validate(i interface{}) error {
	err := dv.validator.Struct(i)

	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return NewBadRequestError(ApiErrorDetails{Message: err.Error()})
		}

		validationErrors := err.(validator.ValidationErrors)

		var apiErrorsDetails []ApiErrorDetails
		for _, validationError := range validationErrors {

			apiErrorDetails := ApiErrorDetails{
				Message: validationError.Translate(dv.translations),
				Param:   validationError.Field(),
			}

			apiErrorsDetails = append(apiErrorsDetails, apiErrorDetails)
		}

		return NewUnprocessableEntityErrors(apiErrorsDetails)
	}

	return nil
}

func (dv *DefaultValidator) UnmarshalAndValidate(bytes []byte, payload any) (any, error) {
	err := json.Unmarshal(bytes, payload)
	if err != nil {
		return nil, err
	}

	err = dv.Validate(payload)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	return payload, nil
}
