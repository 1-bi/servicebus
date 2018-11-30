package runtime

import (
	"github.com/1-bi/servicebus"
	"github.com/1-bi/servicebus/models"
	"github.com/1-bi/servicebus/schema"
	"github.com/1-bi/uerrors"
	"github.com/nats-io/go-nats"
	"github.com/vmihailenco/msgpack"
	"log"
	"time"
)

/**
 * --------------------------------- *
 * package function region
 * --------------------------------- *
 */
func createBaseFuture(mqbrokerUrl string) (bf *baseFuture) {

	bf = new(baseFuture)
	bf.mqbrokerUrl = mqbrokerUrl
	//bf.serviceMng = bsm

	return bf
}

/**
 * --------------------------------- *
 * share object : base future
 * --------------------------------- *
 */
type baseFuture struct {
	nc *nats.Conn
	//serviceMng    *baseServiceAgent
	currentServId string
	mqbrokerUrl   string
	subjectChann  string
	reqMsg        *models.RequestMsg
	timeout       time.Duration
	resultMap     map[string]*models.BaseResult
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
func (this *baseFuture) Await() (coreErr uerrors.CodeError) {

	// --- request message conver to reqmsg ---
	reqMsg := new(schema.ReqMsg)
	reqMsg.Id = this.reqMsg.Id

	paramByteData, err := msgpack.Marshal(this.reqMsg.Params)
	reqMsg.Params = paramByteData

	// --- contruct request message ---
	byteData, err := reqMsg.Marshal(nil)
	if err != nil {
		// --- create new error
		coreErr = uerrors.NewCodeError("E0000015", err.Error())
		log.Fatalf("Error in Request: %v\n", err)
	} else {

		// ---- run sent request ---
		err := this.sentAndReply(this.subjectChann, byteData, this.timeout)
		if err != nil {
			coreErr = uerrors.NewCodeErrorWithPrefix("servbus", "errInSentAndReply", err.Error())
		}
	}

	return coreErr
}

func (this *baseFuture) GetResult() (servicebus.FutureReturnResult, uerrors.CodeError) {

	if this.resultMap == nil {
		return nil, servicebus.Err000002.Build()
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

func (this *baseFuture) convertResultMap(resultItems []*schema.ResultItem) map[string]*models.BaseResult {

	baseResults := make(map[string]*models.BaseResult, 0)
	for _, item := range resultItems {

		resultModel := item.Result

		baseResult := new(models.BaseResult)
		baseResult.ResultRef = resultModel.ResultRef

		if item.Result.Err != nil {

			codeErrModel := resultModel.Err
			codeErr := uerrors.NewCodeErrorWithPrefix(codeErrModel.Prefix, codeErrModel.Code, codeErrModel.MsgBody)

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

func (this *baseFuture) prepareRequest(subjectChann string, reqmsg *models.RequestMsg, timeout time.Duration) {
	this.reqMsg = reqmsg
	this.timeout = timeout
	this.subjectChann = subjectChann
}

/**
 * call and public request
 */
func (this *baseFuture) publishRequest(subjectChann string, reqmsg *models.RequestMsg) error {
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
func (this *baseFuture) newFutureReturnResult(resmap map[string]*models.BaseResult) *baseFutureReturnResult {

	baseFutureRetRes := new(baseFutureReturnResult)

	resErrs := make(map[string]uerrors.CodeError, 0)
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
		baseFutureRetRes.state = servicebus.ALL_ERRORS
	} else if size == len(resRes) {
		baseFutureRetRes.state = servicebus.ALL_COMPLETE
	} else {
		baseFutureRetRes.state = servicebus.ANY_ERRORS
	}

	return baseFutureRetRes
}
