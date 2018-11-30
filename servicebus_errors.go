package servicebus

import "github.com/1-bi/uerrors"

var (
	Err000001 = uerrors.NewCodeErrorWithPrefix(PREFIX, "000000", "Error Test Example")
	// ---- alias error name ----
	Err000002 = uerrors.NewCodeErrorWithPrefix(PREFIX, "000002", "Result map from response is not null.")
	Err000003 = uerrors.NewCodeErrorWithPrefix(PREFIX, "000003", "Object inputted is not a pointer. ")
)
