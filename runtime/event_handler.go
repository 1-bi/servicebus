package runtime

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
)

/**
 * define inner eventhandler
 */
type eventHandlerImpl struct {
}

func (this *eventHandlerImpl) ConvertRequestBody(bingObjFn func() interface{}) {
	returnObj := bingObjFn()
	fmt.Println(returnObj)
}

func (this *eventHandlerImpl) Process(processFn func(bc servicebus.EventbusContext) errors.CodeError) {
	//returnObj := processFn()
	//fmt.Println( returnObj )

	fmt.Println(" message ok  ")
}
