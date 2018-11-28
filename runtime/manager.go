package runtime

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/errors"
	"github.com/1-bi/servicebus/models"
	"github.com/1-bi/servicebus/schema"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"log"
	"sync"
	"time"
)

/**
 * create new base service manager
 */
type baseServiceManager struct {
	baseServiceAgent
}

func NewServiceManager() (servicebus.ServiceManager, error) {

	bsm := new(baseServiceManager)

	bsm.name = "servicebus"

	//mpe := new(MessagePackEncoder)
	//bsm.msgencoder = mpe

	holder := make(map[string][]func(servicebus.ServiceEventHandler), 0)

	bsm.fnHolder = holder

	return bsm, nil
}

/**
 * define base global service handle
 */
func (this *baseServiceManager) On(serviceId string, fn func(servicebus.ServiceEventHandler)) error {

	existedFn := this.fnHolder[serviceId]

	existedFn = append(existedFn, fn)

	// ---- update service function mapping ---

	this.fnHolder[serviceId] = existedFn

	return nil
}

/**
 *
 */
func (this *baseServiceManager) Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (servicebus.Future, errors.CodeError) {

	// --- check object ---

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(nats.DefaultURL)

	// ---- define timeout ----
	f.prepareRequest(this.name, reqmsg, timeout)

	return f, nil
}

func (this *baseServiceManager) FireWithNoReply(serviceId string, runtimeArgs interface{}) errors.CodeError {

	// --- use public handle --

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(nats.DefaultURL)

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
func (this *baseServiceManager) ListenServices() error {

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
func (this *baseServiceManager) doRequest(req *schema.ReqMsg) []*schema.ResultItem {

	servEventhandlers := this.fnHolder[req.Id]

	var wg = new(sync.WaitGroup)

	// ---- process handle ----
	procNum := len(servEventhandlers)

	wg.Add(procNum)

	// --- define result mapping ---
	results := make([]*schema.ResultItem, 0)

	for _, servHandler := range servEventhandlers {

		// --- use goroutine ----
		go func() {

			// --- create handler servert handler implement
			eventHandler := new(eventHandlerImpl)
			eventHandler.serviceManager = this

			// --- predefine handler ----
			servHandler(eventHandler)

			if eventHandler.bindRequestObj != nil {

				// ---- parse object ---

				err := msgpack.Unmarshal(req.Params, eventHandler.bindRequestObj)

				if err != nil {

					codeErr := errors.NewCodeErrorWithPrefix("servbus", "errUnmarshalMessageToBean", "Convert bean error : "+err.Error())

					fmt.Println(codeErr)

				}
			}

			eventHandler.doProcess()

			// ---- get the function name ----
			// ---- conver schema result ----
			eventbusCtx := eventHandler.eventBusCtx
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
			//item.Key = objTypeName
			item.Result = result

			results = append(results, item)

			wg.Done()

		}()

	}

	wg.Wait()

	return results

}

func (this *baseServiceManager) doResponse(serviceId string, resultItems []*schema.ResultItem) *schema.ResMsg {

	// ---- found the result ---

	resMsg := new(schema.ResMsg)
	resMsg.Id = serviceId
	resMsg.Response = resultItems

	return resMsg
}
