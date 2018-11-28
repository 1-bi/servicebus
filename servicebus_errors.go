package servicebus

import "github.com/1-bi/servicebus/errors"

var (
	Err000001 = errors.NewCodeError("000000", "Error Test Example")

	// ---- alias error name ----

	Err000002 = errors.NewCodeErrorWithPrefix("servbus", "000002", "Result map from response is not null.")
	Err000003 = errors.NewCodeErrorWithPrefix("servbus", "000003", "Object inputted is not a pointer. ")
)
