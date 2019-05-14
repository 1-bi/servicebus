package main

import (
	"fmt"
	"github.com/1-bi/servicebus"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/1-bi/servicebus/test"
	"github.com/1-bi/uerrors"
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

	// --- define event handler ---
	serviceManager.On("event.test1", test1Handler)
	serviceManager.On("event.test2", test2Handler)
	serviceManager.On("event.test3", test3Handler)
	serviceManager.On("event.test4", test4Handler)

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}

func test1Handler(handler servicebus.ReqMsgContext) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData string
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data : ")
		fmt.Println(reqData)
		return nil
	})

}

func test2Handler(handler servicebus.ReqMsgContext) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData int
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data2 : ")
		fmt.Println(reqData)
		return nil
	})

}

func test3Handler(handler servicebus.ReqMsgContext) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := make(map[string]string, 0)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data3 : ")
		fmt.Println(reqData)
		return nil
	})
}

func test4Handler(handler servicebus.ReqMsgContext) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := new(test.MockObj1)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data4 : ")
		fmt.Println(reqData)
		return nil
	})
}
