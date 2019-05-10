package servicebus

type Callback interface {
}

type SuccessCallback interface {
	Callback
}

type FailureCallback interface {
}
