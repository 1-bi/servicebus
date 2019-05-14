package test

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/test/fixture"
	"github.com/smartystreets/gunit"
	"testing"
)

// ---- setup method ---
var agent *servicebus.Agent

func TestAgent(t *testing.T) {
	// define method ---
	gunit.Run(new(fixture.AgentFixture), t)

}
