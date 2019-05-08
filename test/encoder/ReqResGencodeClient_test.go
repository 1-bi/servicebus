package encoder

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/encoder"
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/1-bi/uerrors"
	"github.com/nats-io/go-nats"
	"log"
	"testing"
	"time"
)

/**
 * defined publish message
 */
func Test_ReqResGencode_Case1(t *testing.T) {

	var conf *servicebus.AgentConfig
	conf = new(servicebus.AgentConfig)

	conf.SetEncoder((&encoder.GencodeEncoder{}))

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)
	err := agent.SetConfig(conf)

	if err != nil {
		log.Fatal(err)
	}

	var timeout time.Duration
	timeout = 3 * time.Second

	resultRec := func(result servicebus.FutureReturnResult, codeErr uerrors.CodeError) {

		if codeErr != nil {
			log.Panic(codeErr)
			return
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

	agent.FireSyncService("event.req.test1", "hello world , one request ", timeout, resultRec)

}

func Test_Req_Res_Case2(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	resultRec := func(result servicebus.FutureReturnResult, codeErr uerrors.CodeError) {

		if codeErr != nil {
			log.Panic(codeErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

	agent.FireSyncService("event.req.test2", 10001, timeout, resultRec)

}

func Test_Req_Res_Case3(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	baseMap := make(map[string]string, 0)
	baseMap["testkey1"] = "testvalue1"

	resultRec := func(result servicebus.FutureReturnResult, codeErr uerrors.CodeError) {

		if codeErr != nil {
			log.Panic(codeErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

	agent.FireSyncService("event.req.test3", baseMap, timeout, resultRec)

}

/*
func Test_Req_Res_Case4(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	var timeout time.Duration
	timeout = 3 * time.Second

	mockObj1 := new(MockObj1)
	mockObj1.nodeId = "Hello, good boy."
	mockObj1.Age = 20

	resultRec := func(result servicebus.FutureReturnResult, codeErr uerrors.CodeError) {

		if codeErr != nil {
			log.Panic(codeErr)
		}

		var resstr string

		result.ReturnResult(&resstr)

		fmt.Println(" response result ------------------ ")
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

	agent.FireSyncService("event.req.test4", mockObj1, timeout, resultRec)

}
*/
