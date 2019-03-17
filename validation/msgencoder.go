package validation

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

var CurrentMsgEncoderType = byte(0)

// CheckMsgEncoderMatch check the bean validator
func CheckMsgEncoderMatch(fl validator.FieldLevel) bool {

	fmt.Println("content type ")
	fmt.Println(CurrentMsgEncoderType)

	var msgEncoderType byte
	msgEncoderType = CurrentMsgEncoderType

	var headerBytes []byte
	headerBytes = fl.Field().Interface().([]byte)

	if msgEncoderType == headerBytes[0] {
		return true
	} else {
		return false
	}

}
