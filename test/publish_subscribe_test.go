package test

import (
	rt "github.com/1-bi/servicebus/runtime"
	"github.com/nats-io/go-nats"
	"log"
	"testing"
)

/**
 * defined publish message
 */
func Test_Subscribe_Publish_Case1(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	codeErr := agent.FireWithNoReply("event.test1", "hello world")

	if codeErr != nil {
		log.Panic(codeErr)
	}

}

func Test_Subscribe_Publish_Case2(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)
	var a int = 10000

	codeErr := agent.FireWithNoReply("event.test2", a)

	if codeErr != nil {
		log.Panic(codeErr)
	}

}

func Test_Subscribe_Publish_Case3(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	basemap := make(map[string]string, 0)
	basemap["testcase01"] = "testvalue01"

	codeErr := agent.FireWithNoReply("event.test3", basemap)

	if codeErr != nil {
		log.Panic(codeErr)
	}

}

func Test_Subscribe_Publish_Case4(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	mockObj1 := new(MockObj1)
	mockObj1.Name = "Hello, good boy."
	mockObj1.Age = 20

	codeErr := agent.FireWithNoReply("event.test4", mockObj1)

	if codeErr != nil {
		log.Panic(codeErr)
	}

}
