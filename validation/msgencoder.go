package validation

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

func ValidateMsgEncoderType(fl validator.FieldLevel) bool {

	fmt.Println("sdmoo")

	return fl.Field().String() == "awesome"
}
