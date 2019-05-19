package servicebus

import "github.com/1-bi/uerrors"

type embeddedResult struct {
}

func (myself *embeddedResult) Complete(refobj []byte) {

}

func (myself *embeddedResult) Fail(err uerrors.CodeError) {

}

func (myself *embeddedResult) ConvertRepResult() {

}
