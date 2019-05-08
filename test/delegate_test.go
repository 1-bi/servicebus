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

func Test_Delegate_Case1(t *testing.T) {

	// --- create delete object ---

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

func Test_Delegate_Case2(t *testing.T) {

	// --- create delete object ---

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	//var timeout time.Duration
	//timeout = 3 * time.Second

	agent.FireService("event.req.test1", "hello world , one request ")

}

func Test_time(t *testing.T) {

	var tim = time.Now()

	fmt.Println(tim.UnixNano())

}
