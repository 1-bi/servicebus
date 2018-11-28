package models

import (
	"github.com/1-bi/servicebus/schema"
)

/**
 * get response id handle
 */
type ResponseMsg struct {
	Id       string
	Response map[string]*BaseResult
}

func (this *ResponseMsg) ConvertResMsg() *schema.ResMsg {
	return nil
}

func (this *ResponseMsg) ApplyResMsg(reqmsg *schema.ResMsg) {

}
