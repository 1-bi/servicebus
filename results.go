package servicebus

import (
	"github.com/1-bi/servicebus/errors"
	"github.com/vmihailenco/msgpack"
	"log"
)

type BaseResult struct {
	ResultRef []byte
	Err       errors.CodeError
}

func (this *BaseResult) Complete(objref interface{}) {

	bContent, err := msgpack.Marshal(objref)
	if err != nil {
		log.Fatal(err)
	}

	this.ResultRef = bContent
}

func (this *BaseResult) Fail(err errors.CodeError) {
	this.Err = err
}

func (this *BaseResult) IsSuccess() bool {
	if this.Err != nil {
		return false
	}
	return true
}

/**
 * create new Result Object
 */
func NewResult() Result {
	br := new(BaseResult)
	return br
}
