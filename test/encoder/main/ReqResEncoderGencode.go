package main

import (
	"fmt"
	"github.com/1-bi/servicebus"

	"github.com/1-bi/servicebus/encoder"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/1-bi/servicebus/test"
	"github.com/1-bi/uerrors"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"
)

func main() {

	// --- set config envirment ---
	var conf *servicebus.ServerConfig
	conf = new(servicebus.ServerConfig)

	conf.SetEncoder((&encoder.GencodeEncoder{}))

	// ---- create service bus manager ----
	serviceManager, err := rt.NewServiceManager(nats.DefaultURL)
	serviceManager.SetConfig(conf)

	if err != nil {
		log.Fatal(err)
	}

	// --- define event handler ---
	serviceManager.On("event.req.test1", testReqResHandlerEncoder)
	serviceManager.On("event.req.test2", testReqResHandlerEncoder2)
	serviceManager.On("event.req.test3", testReqResHandlerEncoder3)
	serviceManager.On("event.req.test4", testReqResHandlerEncoder4)

	serviceManager.ListenServices()

	// ---- keep program running ----
	runtime.Goexit()

}

func testReqResHandlerEncoder(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData string
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {

		reqData := bc.GetRequestData()

		result := bc.GetResult()
		result.Complete("Ok , I get it .")

		fmt.Println(" request data : ")
		fmt.Println(reqData)

		return nil
	})

}

func testReqResHandlerEncoder2(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData int
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data2 : ")
		fmt.Println(reqData)
		result := bc.GetResult()
		result.Complete("Ok , I get it 20 .")

		return nil
	})

}

func testReqResHandlerEncoder3(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := make(map[string]string, 0)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data3 : ")
		fmt.Println(reqData)

		result := bc.GetResult()
		result.Complete("Ok , I get it 3 .")

		return nil
	})
}

func testReqResHandlerEncoder4(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := new(test.MockObj1)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) uerrors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data4 : ")
		fmt.Println(reqData)

		result := bc.GetResult()
		result.Complete("Ok , I get it 4 .")

		return nil
	})
}
