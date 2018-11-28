package runtime

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
	"github.com/1-bi/servicebus/models"
	"github.com/1-bi/servicebus/schema"
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
}

/**
 * -------------------------------
 *  new service agent
 * -------------------------------
 */
func NewServiceAgent(connectUrl string) servicebus.ServiceAgent {

	bsa := new(baseServiceAgent)
	bsa.natsUrl = connectUrl
	bsa.name = "servicebus"

	holder := make(map[string][]func(servicebus.ServiceEventHandler), 0)

	bsa.fnHolder = holder

	return bsa
}

/**
 * define base global service handle
 */
func (this *baseServiceAgent) On(serviceId string, fn func(servicebus.ServiceEventHandler)) error {

	existedFn := this.fnHolder[serviceId]

	existedFn = append(existedFn, fn)

	// ---- update service function mapping ---

	this.fnHolder[serviceId] = existedFn

	return nil
}

/**
 *
 */
func (this *baseServiceAgent) Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (servicebus.Future, errors.CodeError) {

	// --- check object ---

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(this.natsUrl)

	// ---- define timeout ----
	f.prepareRequest(this.name, reqmsg, timeout)

	return f, nil
}

func (this *baseServiceAgent) FireWithNoReply(serviceId string, runtimeArgs interface{}) errors.CodeError {

	// --- use public handle --

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(this.natsUrl)

	// ---- define timeout ----
	err := f.publishRequest(this.name, reqmsg)

	if err != nil {
		return errors.NewCodeErrorWithPrefix("splider", "0000003000", err.Error())
	}

	return nil

}

/**
 * execute handler for multi veent
 */
func (this *baseServiceAgent) doRequest(req *schema.ReqMsg) []*schema.ResultItem {

	eventFunction := this.fnHolder[req.Id]

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
				eventbusCtx := newEventbusContextImpl(bindRef, this)

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

func (this *baseServiceAgent) doResponse(serviceId string, resultItems []*schema.ResultItem) *schema.ResMsg {

	// ---- found the result ---

	resMsg := new(schema.ResMsg)
	resMsg.Id = serviceId
	resMsg.Response = resultItems

	return resMsg
}
