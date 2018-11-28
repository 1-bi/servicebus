package main

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/1-bi/servicebus/test"
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
	serviceManager.On("event.req.test1", testReqResHandler1)
	//serviceManager.On("event.test2", testReqResHandler2)
	//serviceManager.On("event.test3", testReqResHandler3)
	//serviceManager.On("event.test4", testReqResHandler4)

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}

func testReqResHandler1(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData string
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {

		reqData := bc.GetRequestData()

		result := bc.GetResult()

		result.Complete("Ok , I get it .")

		fmt.Println(" request data : ")
		fmt.Println(reqData)

		return nil
	})

}

func testReqResHandler2(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData int
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data2 : ")
		fmt.Println(reqData)

		return nil
	})

}

func testReqResHandler3(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := make(map[string]string, 0)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data3 : ")
		fmt.Println(reqData)
		return nil
	})
}

func testReqResHandler4(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := new(test.MockObj1)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data4 : ")
		fmt.Println(reqData)
		return nil
	})
}
