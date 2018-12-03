package main

import (
	"github.com/1-bi/servicebus"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/1-bi/servicebus/test"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"
)

func main() {

	// ---- create service bus manager ----

	serviceManager, err := rt.NewServiceManager(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	mockHandler := test.MockHandlerBean1{}

	eventDelagate := servicebus.NewEventsDelegate()

	eventDelagate.AddEvent("event.req.test1", &mockHandler, "EventAction1")
	eventDelagate.AddEvent("event.req.test2", &mockHandler, "EventAction2")
	eventDelagate.AddEvent("event.req.test3", &mockHandler, "EventAction3")
	eventDelagate.AddEvent("event.req.test4", &mockHandler, "EventAction4")

	for key, methodFun := range eventDelagate.AllEventMethods() {
		serviceManager.On(key, methodFun.Interface().(func(servicebus.ServiceEventHandler)))
	}

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}
