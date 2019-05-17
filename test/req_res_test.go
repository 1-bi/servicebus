package test

import (
	"fmt"
	"github.com/1-bi/servicebus"
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
func Test_Req_Res_Case1(t *testing.T) {

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
