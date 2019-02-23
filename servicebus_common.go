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
	PREFIX = "servicebus"

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

// Config set the runtime config for servcie bus
type Config struct {

	// define message conder
	encoder MessageEncoder
}

func (myself *Config) SetEncoder(encoder MessageEncoder) {
	myself.encoder = encoder
}

func (myself *Config) GetEncoder() MessageEncoder {
	return myself.encoder
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
	eventholder       map[string]*EventMapping
	ctlHandlerMapping map[string]map[string]reflect.Value
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

/**
 * get all method by define ---
 */
func (this *EventsDelegate) AllEventMethods() map[string]reflect.Value {

	eventMethods := make(map[string]reflect.Value)

	for serviceId, eventMapping := range this.eventholder {

		structMethods := this.getAllMethodMapByStruct(eventMapping.Ref)

		eventMethods[serviceId] = structMethods[eventMapping.FunctionName]

	}

	return eventMethods

}

/**
 * get method by method name
 */
func (this *EventsDelegate) getAllMethodMapByStruct(structType interface{}) map[string]reflect.Value {

	t := reflect.TypeOf(structType)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	e := t.String()

	mapSize := len(this.ctlHandlerMapping[e])

	if mapSize == 0 {
		// --- create new mapping ---
		this.ctlHandlerMapping[e] = make(map[string]reflect.Value, 0)

		// --- scan all method in this object ----
		objInstRef := reflect.ValueOf(structType)
		typ := objInstRef.Type()

		for i := 0; i < objInstRef.NumMethod(); i++ {

			methodInst := objInstRef.Method(i)
			methodName := typ.Method(i).Name

			this.ctlHandlerMapping[e][methodName] = methodInst
		}

	}

	return this.ctlHandlerMapping[e]

}

func NewEventsDelegate() *EventsDelegate {
	delegate := new(EventsDelegate)
	delegate.eventholder = make(map[string]*EventMapping, 0)
	delegate.ctlHandlerMapping = make(map[string]map[string]reflect.Value, 0)
	return delegate
}
