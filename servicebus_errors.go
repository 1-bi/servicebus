package servicebus

import "github.com/1-bi/servicebus/errors"

var (
	err000001 = errors.NewCodeError("000000", "Error Test Example")

	// ---- alias error name ----

	err000002 = errors.NewCodeErrorWithPrefix("servbus", "000002", "Result map from response is not null.")
	err000003 = errors.NewCodeErrorWithPrefix("servbus", "000003", "Object inputted is not a pointer. ")
)
