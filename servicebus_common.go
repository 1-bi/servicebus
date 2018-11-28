/**
 * =========================================================================== *
 * define all interface in this file
 * =========================================================================== *
 */
package servicebus

import (
	"reflect"
	"time"
)

/**
 * --------------------------------- *
 * constant region
 * --------------------------------- *
 */
const (
	STATUS_DONE    int8 = 5
	STATUS_RUNNING int8 = 1
	STATUS_CANNEL  int8 = 3

	REQ_TIMEOUT = 3 * time.Second

	// ---- value is for  interface "FutureReturnResult"
	NONE         int8 = 0
	ALL_COMPLETE int8 = 1
	ANY_ERRORS   int8 = 2
	ALL_ERRORS   int8 = 3
)

/**
 * --------------------------------- *
 * function region
 * --------------------------------- *
 */
func FunctypeInObject(object interface{}) string {
	objectType := reflect.TypeOf(object)
	function := objectType.Elem().String()
	return function
}

/**

 */
type EventMapping struct {
	EventName string

	// binding reference instance
	Ref interface{}

	/**
	 * define event name
	 */
	FunctionName string
}

/**
 * define event delegate
 */
type EventsDelegate struct {
	eventholder map[string]*EventMapping
}

func (this *EventsDelegate) AddEvent(event string, handler interface{}, functionName string) {

	eventMapping := new(EventMapping)
	eventMapping.EventName = event
	eventMapping.Ref = handler
	eventMapping.FunctionName = functionName

	this.eventholder[event] = eventMapping
}

func (this *EventsDelegate) AllEventsPredefined() map[string]*EventMapping {
	return this.eventholder
}

func NewEventsDelegate() *EventsDelegate {
	delegate := new(EventsDelegate)
	delegate.eventholder = make(map[string]*EventMapping, 0)
	return delegate
}
