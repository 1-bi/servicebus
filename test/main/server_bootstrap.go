package main

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
	rt "github.com/1-bi/servicebus/runtime"
	"log"
	"runtime"
)

func main() {

	// ---- create service bus manager ----

	serviceManager, err := rt.NewServiceManager()

	if err != nil {
		log.Fatal(err)
	}

	// --- define event handler ---
	serviceManager.On("event.test1", test1Handler)
	serviceManager.On("event.test2", test2Handler)

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}

func test1Handler(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		return nil
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		return nil
	})

}

func test2Handler(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		return nil
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		return nil
	})

}
