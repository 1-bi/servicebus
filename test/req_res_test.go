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
		fmt.Println(result)
		fmt.Println(resstr)
		fmt.Println(" response result ------------------ ")
	}

}
