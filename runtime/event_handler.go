package runtime

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
	"log"
	"reflect"
)

/**
 * define inner eventhandler
 */
type eventHandlerImpl struct {
	bindRequestObj interface{}
	bindProcessor  func(bc servicebus.EventbusContext) errors.CodeError
	serviceManager *baseServiceManager
	eventBusCtx    *eventbusContextImpl
}

func (this *eventHandlerImpl) ConvertRequestBody(bingObjFn func() interface{}) {
	this.bindRequestObj = bingObjFn()
}

func (this *eventHandlerImpl) Process(processFn func(bc servicebus.EventbusContext) errors.CodeError) {
	this.bindProcessor = processFn
}

func (this *eventHandlerImpl) doProcess() {

	// request data interface
	var requestData interface{}

	// --- assign  new value object ----
	if reflect.TypeOf(this.bindRequestObj).Kind() == reflect.Ptr {
		requestData = reflect.ValueOf(this.bindRequestObj).Elem()
	} else {
		requestData = this.bindRequestObj
	}

	this.eventBusCtx = newEventbusContextImpl(requestData, this.serviceManager)

	err := this.bindProcessor(this.eventBusCtx)
	if err != nil {
		log.Fatal(err)
	}

}
