package runtime

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/models"
	"github.com/1-bi/servicebus/schema"
	"github.com/1-bi/uerrors"
	"reflect"
	"sync"
	"time"
)

/**
 * build manager
 */
type baseServiceAgent struct {
	/**
	 * 		set current name for mq subject
	 */
	name string

	/**
	 * define function holder
	 */
	fnHolder map[string][]func(servicebus.ServiceEventHandler)

	natsUrl string

	runtimeConf *servicebus.AgentConfig
}

/**
 * -------------------------------
 *  new service agent
 * -------------------------------
 */
//  NewServiceAgent contruct service bus event
func NewServiceAgent(connectUrl string) servicebus.ServiceAgent {

	bsa := new(baseServiceAgent)
	bsa.natsUrl = connectUrl
	bsa.name = "servicebus"

	// --- define defulat conf
	bsa.runtimeConf = new(servicebus.AgentConfig)

	holder := make(map[string][]func(servicebus.ServiceEventHandler), 0)

	bsa.fnHolder = holder

	return bsa
}

func (myself *baseServiceAgent) SetConfig(agentConf *servicebus.AgentConfig) error {

	myself.runtimeConf = agentConf

	// --- ccheck config validate ---

	return nil
}

/**
 * define base global service handle
 */
func (myself *baseServiceAgent) On(serviceId string, fn func(servicebus.ServiceEventHandler)) error {

	existedFn := myself.fnHolder[serviceId]

	existedFn = append(existedFn, fn)

	// ---- update service function mapping ---

	myself.fnHolder[serviceId] = existedFn

	// --- log the message for define ---
	refevent := reflect.ValueOf(fn).String()

	fmt.Println(refevent)

	//log.Printf("Register event :["+ serviceId +"] " + )

	return nil
}

/**
 *
 */
func (myself *baseServiceAgent) Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (servicebus.Future, uerrors.CodeError) {

	// --- check object ---

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(myself.natsUrl)

	// ---- define timeout ----
	f.prepareRequest(myself.name, reqmsg, timeout)

	return f, nil
}

/**
 *
 */
func (myself *baseServiceAgent) FireSyncService(serviceId string, runtimeArgs interface{}, timeout time.Duration, fn func(servicebus.FutureReturnResult, uerrors.CodeError)) {

	// --- check object ---

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(myself.natsUrl)
	f.SetEncoder(myself.runtimeConf.GetEncoder())

	// ---- define timeout ----
	f.prepareRequest(myself.name, reqmsg, timeout)

	codeErr := f.Await()

	if codeErr != nil {
		fn(nil, codeErr)
	}
	result, resErr := f.GetResult()

	fn(result, resErr)

}

/**
 *
 */
func (myself *baseServiceAgent) FireService(serviceId string, runtimeArgs interface{}) error {

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(myself.natsUrl)

	// ---- define timeout ----
	err := f.publishRequest(myself.name, reqmsg)

	if err != nil {
		return uerrors.NewCodeErrorWithPrefix("splider", "0000003000", err.Error())
	}

	return nil

}

/**
 * execute handler for multi veent
 */
func (myself *baseServiceAgent) doRequest(req *schema.ReqMsg) []*schema.ResultItem {

	eventFunction := myself.fnHolder[req.Id]

	var wg = new(sync.WaitGroup)

	// ---- process handle ----
	procNum := len(eventFunction)

	wg.Add(procNum)

	// --- define result mapping ---
	results := make([]*schema.ResultItem, 0)

	for _, servHandler := range eventFunction {

		// --- use goroutine ----
		go func() {
			fmt.Println(servHandler)

			/*

				// --- bind params ------
				bindRef := servHandler.BindParams()

				if bindRef != nil {

					// ---- parse object ---
					err := msgpack.Unmarshal(req.Params, bindRef)

					if err != nil {

						codeErr := errors.NewCodeErrorWithPrefix("servbus", "errUnmarshalMessageToBean", "Convert bean error : "+err.Error())

						fmt.Println(codeErr)

					}
				}

				// --- invoke function with parameter ----
				// ---- get the object type name ----
				objTypeName := servicebus.FunctypeInObject(servHandler)

				// --- create local bus context ---
				eventbusCtx := newEventbusContextImpl(bindRef, myself)

				err := servHandler.Process(eventbusCtx)

				if err != nil {

				}

				// ---- get the function name ----
				// ---- conver schema result ----
				result := new(schema.Result)
				result.ResultRef = eventbusCtx.result.ResultRef

				// --- check the error ---
				if eventbusCtx.result.Err != nil {

					error := new(schema.CodeError)
					error.Code = eventbusCtx.result.Err.Code()
					error.Prefix = eventbusCtx.result.Err.Prefix()
					error.MsgBody = eventbusCtx.result.Err.MsgBody()

					result.Err = error

				}

				item := new(schema.ResultItem)
				item.Key = objTypeName
				item.Result = result

				results = append(results, item)
			*/
			wg.Done()

		}()

	}

	wg.Wait()

	return results

}

func (myself *baseServiceAgent) doResponse(serviceId string, resultItems []*schema.ResultItem) *schema.ResMsg {

	// ---- found the result ---

	resMsg := new(schema.ResMsg)
	resMsg.Id = serviceId
	resMsg.Response = resultItems

	return resMsg
}
