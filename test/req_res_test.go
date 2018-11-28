package test

import (
	"fmt"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/nats-io/go-nats"
	"log"
	"testing"
	"time"
)

/**
 * defined publish message
 */
func Test_Req_Res_Case1(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	f, codeErr := agent.Fire("event.req.test1", "hello world , one request ", timeout)

	if codeErr != nil {
		log.Panic(codeErr)
	} else {

		codeErr = f.Await()

		if codeErr != nil {

			log.Panic(codeErr)

		}

		result, resErr := f.GetResult()

		if resErr != nil {
			log.Panic(resErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

}

func Test_Req_Res_Case2(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	f, codeErr := agent.Fire("event.req.test2", 10001, timeout)

	if codeErr != nil {
		log.Panic(codeErr)
	} else {

		codeErr = f.Await()

		if codeErr != nil {

			log.Panic(codeErr)

		}

		result, resErr := f.GetResult()

		if resErr != nil {
			log.Panic(resErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}
}

func Test_Req_Res_Case3(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	baseMap := make(map[string]string, 0)
	baseMap["testkey1"] = "testvalue1"

	f, codeErr := agent.Fire("event.req.test3", baseMap, timeout)

	if codeErr != nil {
		log.Panic(codeErr)
	} else {

		codeErr = f.Await()

		if codeErr != nil {

			log.Panic(codeErr)

		}

		result, resErr := f.GetResult()

		if resErr != nil {
			log.Panic(resErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}
}

func Test_Req_Res_Case4(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	mockObj1 := new(MockObj1)
	mockObj1.Name = "Hello, good boy."
	mockObj1.Age = 20

	f, codeErr := agent.Fire("event.req.test4", mockObj1, timeout)

	if codeErr != nil {
		log.Panic(codeErr)
	} else {

		codeErr = f.Await()

		if codeErr != nil {

			log.Panic(codeErr)

		}

		result, resErr := f.GetResult()

		if resErr != nil {
			log.Panic(resErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}
}
