package utils

import "github.com/go-playground/validator/v10"

type errorResponse struct {
	Field string
	Tag   string
	Value string
}

func (e errorResponse) Error() string {
	return e.Field + ":\"" + e.Tag + "=" + e.Value
}

func ValidateStruct(in any) []*errorResponse {
	var errors []*errorResponse
	validate := validator.New()

	err := validate.Struct(in)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}
