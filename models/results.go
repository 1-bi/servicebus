package models

import (
	"github.com/1-bi/uerrors"
	"github.com/vmihailenco/msgpack"
	"log"
)

type BaseResult struct {
	ResultRef []byte
	Err       uerrors.CodeError
}

func (this *BaseResult) Complete(objref interface{}) {

	bContent, err := msgpack.Marshal(objref)
	if err != nil {
		log.Fatal(err)
	}

	this.ResultRef = bContent
}

func (this *BaseResult) Fail(err uerrors.CodeError) {
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
func NewResult() *BaseResult {
	br := new(BaseResult)
	return br
}
