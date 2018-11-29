package test

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
)

/**
 * defined mock handler bean
 */
type MockHandlerBean1 struct {
}

func (this *MockHandlerBean1) EventAction1(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData string
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {

		reqData := bc.GetRequestData()

		result := bc.GetResult()
		result.Complete("Ok , I get it in MockHandlerBean1 .")

		fmt.Println(" request data : ")
		fmt.Println(reqData)

		return nil
	})
}

func (this *MockHandlerBean1) EventAction2(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		var reqData int
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data2 : ")
		fmt.Println(reqData)
		result := bc.GetResult()
		result.Complete("Ok , I get it twice in MockHandlerBean1  .")

		return nil
	})
}

func (this *MockHandlerBean1) EventAction3(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := make(map[string]string, 0)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data3 : ")
		fmt.Println(reqData)

		result := bc.GetResult()
		result.Complete("Ok , I get it third in MockHandlerBean1 .")

		return nil
	})

}

func (this *MockHandlerBean1) EventAction4(handler servicebus.ServiceEventHandler) {

	handler.ConvertRequestBody(func() interface{} {
		reqData := new(MockObj1)
		return &reqData
	})

	handler.Process(func(bc servicebus.EventbusContext) errors.CodeError {
		reqData := bc.GetRequestData()

		fmt.Println(" request data4 : ")
		fmt.Println(reqData)

		result := bc.GetResult()
		result.Complete("Ok , I get it forth in MockHandlerBean1 .")

		return nil
	})
}
