package test

import (
	"github.com/1-bi/servicebus/test/fixture"
	"github.com/smartystreets/gunit"
	"testing"
)

// ---- setup method ---

func TestEtcdServiceOperations(t *testing.T) {
	// define method ---
	gunit.Run(new(fixture.EtcdServiceOperationsFixture), t)
}
