/**
 * =========================================================================== *
 * define all interface in this file
 * =========================================================================== *
 */
package servicebus

import (
	"github.com/1-bi/servicebus/errors"
	"time"
)

/**
 * create all interface context
 */

type ServiceEvent interface {

	/**
	 * call one service handle
	 */
	On(serviceId string, fn ServiceEventHandler) error

	/**
	 * create fire service event object , but the object is not runt
	 * It could be run when the future is called by await
	 * default is 100 milion seonc
	 */
	Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (Future, errors.CodeError)

	/**
	 * create fire servic event object without no reply after handling
	 */
	FireWithNoReply(serviceId string, runtimeArgs interface{}) errors.CodeError
}

/**
 * create the service agent
 */
type ServiceAgent interface {
	ServiceEvent
}

/**
 * contruct service manager
 */
type ServiceManager interface {
	ServiceEvent

	/**
	 * boot listen service
	 */
	ListenServices() error
}

// --- create bus context ----
type EventbusContext interface {
	GetRoot() interface{}

	/**
	 * Get the result inteface
	 */
	GetResult() Result

	/**
	 * support service event
	 */
	GetServiceEvent() ServiceEvent
}

/**
 * define base result
 */
type Result interface {

	/**
	 *	complement object
	 */
	Complete(refobj interface{})

	/**
	 *  add the result
	 */
	Fail(err errors.CodeError)
}

/**
 * contruct service event handler
 * @deplecated this object is not to use
 */
type ServiceEventHandler interface {
	/**
	 * define root object
	 */
	BindParams() interface{}

	Process(bc EventbusContext) errors.CodeError
}

/**
 * return the service event futrue
 */
type Future interface {

	/**
	 * wait the event complete
	 */
	Await() errors.CodeError

	/**
	 * get the result object after await
	 */
	GetResult() (FutureReturnResult, errors.CodeError)
}

/**
 *  This interface will map the futurn result handle
 */
type FutureReturnResult interface {

	/**
	 * define,  ALL_COMPLETE ,  ANY_ERRORS , ALL_ERRORS
	 */
	State() int8

	/**
	 *  return all error from service event running
	 */
	Errors(procName string) errors.CodeError

	/**
	 * rturn all return Results
	 */
	ReturnResults(procName string, inReturn interface{}) errors.CodeError

	/**
	 * get the first error directly
	 */
	Error() errors.CodeError

	/**
	 * get the first result directly
	 */
	ReturnResult(inReturn interface{}) errors.CodeError
}

/**
 * define root type decoder
 */
type RootTypeDecoder interface {
	ProvideRootRef() interface{}
}
