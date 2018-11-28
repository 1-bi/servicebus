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
func Test_Case1_Publish(t *testing.T) {

	natsUrl := nats.DefaultURL

	agent := rt.NewServiceAgent(natsUrl)

	codeErr := agent.FireWithNoReply("event.test1", "hello world")

	if codeErr != nil {
		log.Panic(codeErr)
	}

}
