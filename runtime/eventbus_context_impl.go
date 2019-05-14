package runtime

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/models"
)

/**
 * --------------------------------- *
 * share object : base bus
 * --------------------------------- *
 */
type eventbusContextImpl struct {
	params       interface{}
	result       *models.BaseResult
	serviceEvent servicebus.ServiceEvent
}

func (this *eventbusContextImpl) setParams(in interface{}) {
	this.params = in
}

func (this *eventbusContextImpl) GetRequestData() interface{} {
	return this.params
}

func (this *eventbusContextImpl) GetResult() servicebus.Result {
	return nil
}

func (this *eventbusContextImpl) GetServiceEvent() servicebus.ServiceEvent {
	// --- define fire object ---
	return this.serviceEvent
}

func newEventbusContextImpl(params interface{}, serviceEvent servicebus.ServiceEvent) *eventbusContextImpl {
	eci := new(eventbusContextImpl)
	eci.result = new(models.BaseResult)
	eci.params = params
	eci.serviceEvent = serviceEvent
	return eci
}
