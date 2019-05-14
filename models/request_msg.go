package models

import "github.com/1-bi/servicebus/schema"

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
	msg    *schema.ReqMsg
}

func NewRequestMsg() *RequestMsg {
	msg := new(RequestMsg)
	return msg
}

func ConvertRequestMsgFromReqMsg(reqMsg *schema.ReqMsg) *RequestMsg {
	msg := new(RequestMsg)
	msg.msg = reqMsg
	return msg
}

func (this *RequestMsg) ConvertReqMsg() *schema.ReqMsg {
	return nil
}

func (this *RequestMsg) ApplyReqMsg(reqmsg *schema.ReqMsg) {

}
