package errors

import (
	"bytes"
	"fmt"
)

type CodeError interface {
	error
	Code() string
	Prefix() string
	MsgBody() string
	/**
	 * build err object
	 */
	Build(args ...string) CodeError
}

type ValidatorError interface {
	CodeError
}

/**
 * build code error
 */
type baseCodeError struct {
	code    string
	msgBody string
	prefix  string
}

/**
 * implment error function
 */
func (this *baseCodeError) Error() string {

	// --- check message content ----
	buf := bytes.NewBufferString("[")

	if this.prefix != "" {
		buf.WriteString(this.prefix)
		buf.WriteString(":")
	}

	buf.WriteString(this.code)
	buf.WriteString("]: ")

	buf.WriteString(this.msgBody)
	buf.WriteString("\n ")

	return fmt.Sprintf(buf.String())
}

func (this *baseCodeError) Code() string {
	return this.code
}

func (this *baseCodeError) Prefix() string {
	return this.prefix
}

func (this *baseCodeError) MsgBody() string {
	return this.msgBody
}

func (this *baseCodeError) Build(args ...string) CodeError {
	return this
}

/**
 *   create public method code error , use simple code handle --
 */
func NewCodeError(code string, msgContent ...string) CodeError {

	ce := baseCodeError{code, "", ""}

	// --- check message content ----
	buf := bytes.NewBufferString("")

	for _, msg := range msgContent {
		buf.WriteString(msg)
	}

	ce.msgBody = buf.String()

	return &ce
}

/**
 * create new code error
 */
func NewCodeErrorWithPrefix(prefix string, code string, msgContent ...string) CodeError {

	ce := baseCodeError{code, "", prefix}

	// --- check message content ----
	buf := bytes.NewBufferString("")

	for _, msg := range msgContent {
		buf.WriteString(msg)
	}

	ce.msgBody = buf.String()

	return &ce
}
