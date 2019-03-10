package runtime

import (
	"fmt"
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/models"
	"github.com/1-bi/servicebus/schema"
	"github.com/1-bi/servicebus/validation"
	"github.com/1-bi/uerrors"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

/**
 * create new base service manager
 */
type baseServiceManager struct {
	baseServiceAgent
	natsUrl    string
	msgEncoder servicebus.MessageEncoder
	validate   *validator.Validate
}

func NewServiceManager(natsUrl string) (servicebus.ServiceManager, error) {

	bsm := new(baseServiceManager)

	bsm.name = "servicebus"

	//mpe := new(MessagePackEncoder)
	//bsm.msgencoder = mpe

	holder := make(map[string][]func(servicebus.ServiceEventHandler), 0)
	bsm.fnHolder = holder
	bsm.natsUrl = natsUrl
	bsm.validate = validator.New()

	return bsm, nil
}

// SetConfig init the servicebus instance
func (myself *baseServiceManager) SetConfig(conf *servicebus.ServerConfig) error {

	// --- check the config envirment ----

	err := myself.validateMessageEncoder(conf.GetEncoder())
	if err != nil {
		return err
	}
	myself.msgEncoder = conf.GetEncoder()

	return nil
}

func (myself *baseServiceManager) validateMessageEncoder(encoder servicebus.MessageEncoder) error {

	type MessageEncodertype struct {
		encoderInst interface{} `validate:"validate-msgencoder"`
	}

	met := MessageEncodertype{}
	myself.validate.RegisterValidation("validate-msgencoder", validation.ValidateMsgEncoderMatch)

	err := myself.validate.Struct(met)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}

	return nil

}

/**
 * define base global service handle
 */
func (myself *baseServiceManager) On(serviceId string, fn func(servicebus.ServiceEventHandler)) error {

	existedFn := myself.fnHolder[serviceId]

	existedFn = append(existedFn, fn)

	// ---- update service function mapping ---
	myself.fnHolder[serviceId] = existedFn

	refevent := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	msg := strings.Join([]string{"Register event [", serviceId, "], function name [", refevent, "]."}, "")
	log.Println(msg)
	return nil
}

/**
 *
 */
func (myself *baseServiceManager) Fire(serviceId string, runtimeArgs interface{}, timeout time.Duration) (servicebus.Future, uerrors.CodeError) {

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

func (myself *baseServiceManager) FireWithNoReply(serviceId string, runtimeArgs interface{}) uerrors.CodeError {

	// --- use public handle --

	// ---- create request msg ----
	reqmsg := models.NewRequestMsg()
	reqmsg.Id = serviceId
	reqmsg.Params = runtimeArgs

	// ---- create current event ---
	f := createBaseFuture(nats.DefaultURL)

	// ---- define timeout ----
	err := f.publishRequest(myself.name, reqmsg)

	if err != nil {
		return uerrors.NewCodeErrorWithPrefix("splider", "0000003000", err.Error())
	}

	return nil

}

func (myself *baseServiceManager) initValidation() error {

	myself.validate.RegisterValidation("check-encoder-match", validation.ValidateMsgEncoderMatch)

	return nil
}

/**
 * start up boot and listen service
 * Listen service use default request / rply mode
 */
func (myself *baseServiceManager) ListenServices() error {

	nc, err := nats.Connect(myself.natsUrl)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	err = myself.initValidation()
	if err != nil {
		return err
	}

	//queue := "default"

	subj := myself.name

	_, suberr := nc.Subscribe(subj, func(msg *nats.Msg) {

		maxLength := len(msg.Data)
		// --- check the conder and encoder --
		headerBytes := msg.Data[:8]
		bodyBytes := msg.Data[8 : maxLength-1]

		// --- check and validate data ---
		msgVal := new(validation.MsgCandidateVali)
		msgVal.Header = headerBytes

		fmt.Println("recieve size  ")
		fmt.Println(maxLength)

		fmt.Println(headerBytes)
		fmt.Println(bodyBytes)

		// --- get the header flag ---
		// ---- convert to req message ---
		reqMsg := new(schema.ReqMsg)
		reqMsg.Unmarshal(bodyBytes)

		// --- get service process by service id ----
		resmap := myself.doRequest(reqMsg)

		// --- result map ---
		resMsg := myself.doResponse(reqMsg.Id, resmap)

		byteContent, err := resMsg.Marshal(nil)

		if err != nil {
			log.Fatal(err)
		}

		// ---- reply message
		nc.Publish(msg.Reply, byteContent)

	})
	nc.Flush()

	if suberr != nil {

	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)

	return nil
}

/**
 * execute handler for multi veent
 */
func (myself *baseServiceManager) doRequest(req *schema.ReqMsg) []*schema.ResultItem {

	servEventhandlers := myself.fnHolder[req.Id]

	var wg = new(sync.WaitGroup)

	// ---- process handle ----
	procNum := len(servEventhandlers)

	wg.Add(procNum)

	// --- define result mapping ---
	results := make([]*schema.ResultItem, 0)

	for _, servHandler := range servEventhandlers {

		// --- use goroutine ----
		go func() {

			funName := runtime.FuncForPC(reflect.ValueOf(servHandler).Pointer()).Name()

			// --- create handler servert handler implement
			eventHandler := new(eventHandlerImpl)
			eventHandler.serviceManager = myself

			// --- predefine handler ----
			servHandler(eventHandler)

			if eventHandler.bindRequestObj != nil {

				// ---- parse object ---

				err := msgpack.Unmarshal(req.Params, eventHandler.bindRequestObj)

				if err != nil {

					codeErr := uerrors.NewCodeErrorWithPrefix("servbus", "errUnmarshalMessageToBean", "Convert bean error : "+err.Error())

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
			item.Key = funName
			item.Result = result

			results = append(results, item)

			wg.Done()

		}()

	}

	wg.Wait()

	return results

}

func (myself *baseServiceManager) doResponse(serviceId string, resultItems []*schema.ResultItem) *schema.ResMsg {

	// ---- found the result ---

	resMsg := new(schema.ResMsg)
	resMsg.Id = serviceId
	resMsg.Response = resultItems

	return resMsg
}
