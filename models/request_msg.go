package models

/**
 * --------------------------------- *
 * fixture region
 * --------------------------------- *
 */
/**
 * service stub convert message
 */
type RequestMsg struct {
	Id     string
	Params interface{}
}

func NewRequestMsg() *RequestMsg {
	msg := new(RequestMsg)
	return msg
}
