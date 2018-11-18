/**
 * =========================================================================== *
 * define all interface in this file
 * =========================================================================== *
 */
package servicebus

import (
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"gitlab.com/vicenteyuen/tmall-splider/errors"
	"gitlab.com/vicenteyuen/tmall-splider/servicebus/schema"
	"log"
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
 * struct region
 * --------------------------------- *
 */
/**
 * service stub convert message
 */
type RequestMsg struct {
	Id     string
	Params interface{}
}

/**
 * get response id handle
 */
type ResponseMsg struct {
	Id       string
	Response map[string]*BaseResult
}

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
 * create new message
 */
func NewRequestMsg(id string, params interface{}) (reqmsg *RequestMsg) {
	reqmsg = new(RequestMsg)
	reqmsg.Id = id
	reqmsg.Params = params

	// ---- get the type mapping
	return reqmsg
}

/**
 * --------------------------------- *
 * package function region
 * --------------------------------- *
 */
func createBaseFuture(bsm *baseServiceAgent, mqbrokerUrl string) (bf *baseFuture) {

	bf = new(baseFuture)
	bf.mqbrokerUrl = mqbrokerUrl
	bf.serviceMng = bsm

	return bf
}

/**
 * --------------------------------- *
 * share object : base bus
 * --------------------------------- *
 */
type eventbusContextImpl struct {
	params       interface{}
	result       *BaseResult
	serviceEvent ServiceEvent
}

func (this *eventbusContextImpl) setParams(in interface{}) {
	this.params = in
}

func (this *eventbusContextImpl) GetRoot() interface{} {
	return this.params
}

func (this *eventbusContextImpl) GetResult() Result {
	return this.result
}

func (this *eventbusContextImpl) GetServiceEvent() ServiceEvent {
	// --- define fire object ---
	return this.serviceEvent
}

func newEventbusContextImpl(params interface{}, serviceEvent ServiceEvent) *eventbusContextImpl {
	eci := new(eventbusContextImpl)
	eci.result = new(BaseResult)
	eci.params = params
	eci.serviceEvent = serviceEvent
	return eci
}

/**
 * --------------------------------- *
 * share object : base future
 * --------------------------------- *
 */
type baseFuture struct {
	nc            *nats.Conn
	serviceMng    *baseServiceAgent
	currentServId string
	mqbrokerUrl   string
	subjectChann  string
	reqMsg        *RequestMsg
	timeout       time.Duration
	resultMap     map[string]*BaseResult
}

/**
 * return current future status
 */
func (this *baseFuture) GetStatus() int8 {
	return 0
}

/**
 * sent the message to mq in this function
 */
func (this *baseFuture) Await() (coreErr errors.CodeError) {

	// --- request message conver to reqmsg ---
	reqMsg := new(schema.ReqMsg)
	reqMsg.Id = this.reqMsg.Id

	paramByteData, err := msgpack.Marshal(this.reqMsg.Params)
	reqMsg.Params = paramByteData

	// --- contruct request message ---
	byteData, err := reqMsg.Marshal(nil)
	if err != nil {
		// --- create new error
		coreErr = errors.NewCodeError("E0000015", err.Error())
		log.Fatalf("Error in Request: %v\n", err)
	} else {

		// ---- run sent request ---
		err := this.sentAndReply(this.subjectChann, byteData, this.timeout)
		if err != nil {
			coreErr = errors.NewCodeErrorWithPrefix("servbus", "errInSentAndReply", err.Error())
		}
	}

	return coreErr
}

func (this *baseFuture) GetResult() (FutureReturnResult, errors.CodeError) {

	if this.resultMap == nil {
		return nil, err000002
	}

	// ---- check data error
	baseFuReResult := this.newFutureReturnResult(this.resultMap)

	return baseFuReResult, nil
}

/**
 * sent subject request
 */
func (this *baseFuture) sentAndReply(subject string, content []byte, timeout time.Duration) error {

	nc, err := nats.Connect(this.mqbrokerUrl)
	if err != nil {
		return err
	}
	defer nc.Close()

	msg, err := nc.Request(subject, content, timeout)
	if err != nil {
		if nc.LastError() != nil {
			return nc.LastError()
		}
		return err
	}
	// --- convert to replied object ---
	resMsg := new(schema.ResMsg)
	resMsg.Unmarshal(msg.Data)

	// --- create response mapping ------
	resultItems := this.doReply(resMsg)
	// --- convert result result  ---
	this.resultMap = this.convertResultMap(resultItems)

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subject, content)
	}

	return nil
}

func (this *baseFuture) convertResultMap(resultItems []*schema.ResultItem) map[string]*BaseResult {

	baseResults := make(map[string]*BaseResult, 0)
	for _, item := range resultItems {

		resultModel := item.Result

		baseResult := new(BaseResult)
		baseResult.ResultRef = resultModel.ResultRef

		if item.Result.Err != nil {

			codeErrModel := resultModel.Err
			codeErr := errors.NewCodeErrorWithPrefix(codeErrModel.Prefix, codeErrModel.Code, codeErrModel.MsgBody)

			baseResult.Fail(codeErr)

		}
		baseResults[item.Key] = baseResult
	}

	return baseResults
}

func (this *baseFuture) send(subject string, content []byte) error {
	nc, err := nats.Connect(this.mqbrokerUrl)
	if err != nil {
		return err
	}
	defer nc.Close()

	err = nc.Publish(subject, content)
	if err != nil {
		if nc.LastError() != nil {
			return nc.LastError()
		}
		return err
	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subject, content)
	}

	return nil
}

/**
 * receive message from response
 */
func (this *baseFuture) doReply(resMsg *schema.ResMsg) []*schema.ResultItem {

	return resMsg.Response

}

func (this *baseFuture) prepareRequest(subjectChann string, reqmsg *RequestMsg, timeout time.Duration) {
	this.reqMsg = reqmsg
	this.timeout = timeout
	this.subjectChann = subjectChann
}

/**
 * call and public request
 */
func (this *baseFuture) publishRequest(subjectChann string, reqmsg *RequestMsg) error {
	this.reqMsg = reqmsg
	this.subjectChann = subjectChann

	reqMsg := new(schema.ReqMsg)
	reqMsg.Id = this.reqMsg.Id

	paramByteData, err := msgpack.Marshal(this.reqMsg.Params)
	reqMsg.Params = paramByteData
	if err != nil {
		return err
	}

	// --- contruct request message ---
	byteData, err := reqMsg.Marshal(nil)

	this.send(this.subjectChann, byteData)

	return nil
}

/**
 * receive response message
 */
func (this *baseFuture) newFutureReturnResult(resmap map[string]*BaseResult) *baseFutureReturnResult {

	baseFutureRetRes := new(baseFutureReturnResult)

	resErrs := make(map[string]errors.CodeError, 0)
	resRes := make(map[string][]byte, 0)

	size := len(resmap)

	for procName, resultObj := range resmap {

		if resultObj.Err != nil {
			resErrs[procName] = resultObj.Err
		} else {

			// --- add mapping ---
			resRes[procName] = resultObj.ResultRef
		}

	}
	baseFutureRetRes.resErrs = resErrs
	baseFutureRetRes.resRes = resRes

	// --- state ----
	if size == len(resErrs) {
		baseFutureRetRes.state = ALL_ERRORS
	} else if size == len(resRes) {
		baseFutureRetRes.state = ALL_COMPLETE
	} else {
		baseFutureRetRes.state = ANY_ERRORS
	}

	return baseFutureRetRes
}

/**
 * --------------------------------- *
 * share object : base baseFutureReturnResult
 * --------------------------------- *
 */
type baseFutureReturnResult struct {
	state   int8
	resErrs map[string]errors.CodeError
	resRes  map[string][]byte
}

/**
 * define,  ALL_COMPLETE ,  ANY_ERRORS , ALL_ERRORS
 */
func (this *baseFutureReturnResult) State() int8 {
	return this.state
}

/**
 *  return all error from service event running
 */
func (this *baseFutureReturnResult) Errors(procName string) errors.CodeError {
	return this.resErrs[procName]
}

/**
 * rturn all return Results
 */
func (this *baseFutureReturnResult) ReturnResults(procName string, inReturn interface{}) errors.CodeError {

	// ---- validate type ---
	if reflect.TypeOf(inReturn).Kind() != reflect.Ptr {
		return err000003
	}

	resobj := this.resRes[procName]
	msgpack.Unmarshal(resobj, inReturn)
	return nil

}

/**
 * get the first error directly
 */
func (this *baseFutureReturnResult) Error() errors.CodeError {
	// --- get the first error --
	for _, errObj := range this.resErrs {
		return errObj
	}
	return nil
}

/**
 * get the first result directly
 */
func (this *baseFutureReturnResult) ReturnResult(inReturn interface{}) errors.CodeError {

	// ---- validate type ---
	if reflect.TypeOf(inReturn).Kind() != reflect.Ptr {
		return err000003
	}

	var objRef []byte
	for _, tmpObjRef := range this.resRes {
		objRef = tmpObjRef
		break
	}

	msgpack.Unmarshal(objRef, inReturn)

	return nil
}

/**
 * define event delegate
 */
type EventsDelegate struct {
	eventholder map[string]ServiceeventHandler
}

func (this *EventsDelegate) AddEvent(event string, handler ServiceeventHandler) {
	this.eventholder[event] = handler
}

func (this *EventsDelegate) AllEventsPredefined() map[string]ServiceeventHandler {
	return this.eventholder
}

func NewEventsDelegate() *EventsDelegate {
	delegate := new(EventsDelegate)
	delegate.eventholder = make(map[string]ServiceeventHandler, 0)
	return delegate
}
