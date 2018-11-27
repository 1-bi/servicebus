package servicebus

import (
	"fmt"
	"github.com/1-bi/servicebus/errors"
	"github.com/1-bi/servicebus/schema"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"log"
	"sync"
	"time"
)

func NewServiceManager() (ServiceManager, error) {

	bsm := new(baseServiceManager)

	bsm.name = "servicebus"

	//mpe := new(MessagePackEncoder)
	//bsm.msgencoder = mpe

	holder := make(map[string][]ServiceEventHandler, 0)

	bsm.fnHolder = holder

	return bsm, nil
}

/**
 * define base global service handle
 */
func (this *baseServiceAgent) On(serviceId string, fn ServiceEventHandler) error {

	existedFn := this.fnHolder[serviceId]

	existedFn = append(existedFn, fn)

	// ---- update service function mapping ---

	this.fnHolder[serviceId] = existedFn

	return nil
}

/**
 *
 */
func (this *baseServiceAgent) Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (Future, errors.CodeError) {

	// --- check object ---

	// ---- create request msg ----
	reqmsg := NewRequestMsg(serviceId, runtimeArgs)

	// ---- create current event ---
	f := createBaseFuture(this, nats.DefaultURL)

	// ---- define timeout ----
	f.prepareRequest(this.name, reqmsg, timeout)

	return f, nil
}

func (this *baseServiceAgent) FireWithNoReply(serviceId string, runtimeArgs interface{}) errors.CodeError {

	// --- use public handle --

	// ---- create request msg ----
	reqmsg := NewRequestMsg(serviceId, runtimeArgs)

	// ---- create current event ---
	f := createBaseFuture(this, nats.DefaultURL)

	// ---- define timeout ----
	err := f.publishRequest(this.name, reqmsg)

	if err != nil {
		return errors.NewCodeErrorWithPrefix("splider", "0000003000", err.Error())
	}

	return nil

}

/**
 * start up boot and listen service
 * Listen service use default request / rply mode
 */
func (this *baseServiceAgent) ListenServices() error {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	//queue := "default"

	subj := this.name

	nc.Subscribe(subj, func(msg *nats.Msg) {

		// ---- convert to req message ---
		reqMsg := new(schema.ReqMsg)

		reqMsg.Unmarshal(msg.Data)
		// --- get service process by service id ----
		resmap := this.doRequest(reqMsg)

		// --- result map ---
		resMsg := this.doResponse(reqMsg.Id, resmap)

		byteContent, err := resMsg.Marshal(nil)

		if err != nil {
			log.Fatal(err)
		}

		// ---- reply message
		nc.Publish(msg.Reply, byteContent)

	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)

	return nil
}

/**
 * execute handler for multi veent
 */
func (this *baseServiceAgent) doRequest(req *schema.ReqMsg) []*schema.ResultItem {

	serveventHandler := this.fnHolder[req.Id]

	var wg = new(sync.WaitGroup)

	// ---- process handle ----
	procNum := len(serveventHandler)

	wg.Add(procNum)

	// --- define result mapping ---
	results := make([]*schema.ResultItem, 0)

	for _, servHandler := range serveventHandler {

		// --- use goroutine ----
		go func() {

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
			objTypeName := FunctypeInObject(servHandler)

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
