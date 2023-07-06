package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

// Use a single instance of Validator for all apis.
//
// Validate is designed to be thread-safe and used as a singleton instance.
var validate *validator.Validate = validator.New()

// A universal translator for validation messages.
var trans ut.Translator

func init() {
	// Registers a function to use tag names for StructFields, in validation messages.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, trans)
}

// translateError translate validation errors and returns a list of errors.
func translateError(err error) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr)
	}
	return errs
}

// RunValidation runs the struct validation on given data and returns the potential error list.
func RunValidation(data interface{}) []error {
	errs := []error{}
	if err := validate.Struct(data); err != nil {
		errs = translateError(err)
	}
	return errs
}
