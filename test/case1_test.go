package test

import (
	"github.com/1-bi/servicebus"
	"log"
	"testing"
)

/**
 * defined publish message
 */
func Test_Case1_Publish(t *testing.T) {

	agent := servicebus.NewServiceAgent()

	codeErr := agent.FireWithNoReply("event.test", "hello world")

	if codeErr != nil {
		log.Panic(codeErr)
	}

}
