package runtime

import (
	"bytes"
	"fmt"
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
	encoder       servicebus.MessageEncoder
}

func (myself *baseFuture) SetEncoder(encoder servicebus.MessageEncoder) {
	myself.encoder = encoder
}

/**
 * return current future status
 */
func (myself *baseFuture) GetStatus() int8 {
	return 0
}

func (myself *baseFuture) contructMsgByte(body []byte) []byte {

	header := make([]byte, 8)

	// encoder type ---
	header[0] = myself.encoder.GetType()

	// --- is empty ---
	header[1] = 0
	header[2] = 0
	header[3] = 0
	header[4] = 0
	header[5] = 0
	header[6] = 0
	header[7] = 0

	msgBytes := [][]byte{header, body}
	msgContentBytes := bytes.Join(msgBytes, []byte{})

	return msgContentBytes

}

/**
 * sent the message to mq in this function
 */
func (myself *baseFuture) Await() (coreErr uerrors.CodeError) {

	// --- request message conver to reqmsg ---
	reqMsg := new(schema.ReqMsg)
	reqMsg.Id = myself.reqMsg.Id

	paramByteData, err := msgpack.Marshal(myself.reqMsg.Params)
	reqMsg.Params = paramByteData

	// --- contruct request message ---
	byteData, err := reqMsg.Marshal(nil)
	if err != nil {
		// --- create new error
		coreErr = uerrors.NewCodeError("E0000015", err.Error())
		log.Fatalf("Error in Request: %v\n", err)
		return coreErr
	}

	// --- contruct message content ---
	byteData = myself.contructMsgByte(byteData)

	sendSize := len(byteData)
	fmt.Println("send size : ")
	fmt.Println(sendSize)

	// ---- sent object ---
	err = myself.sentAndReply(myself.subjectChann, byteData, myself.timeout)
	if err != nil {
		coreErr = uerrors.NewCodeErrorWithPrefix("servbus", "errInSentAndReply", err.Error())
	}

	return coreErr
}

func (myself *baseFuture) GetResult() (servicebus.FutureReturnResult, uerrors.CodeError) {

	if myself.resultMap == nil {
		return nil, servicebus.Err000002.Build()
	}

	// ---- check data error
	baseFuReResult := myself.newFutureReturnResult(myself.resultMap)

	return baseFuReResult, nil
}

/**
 * sent subject request
 */
func (myself *baseFuture) sentAndReply(subject string, content []byte, timeout time.Duration) error {

	nc, err := nats.Connect(myself.mqbrokerUrl)
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
	resultItems := myself.doReply(resMsg)
	// --- convert result result  ---
	myself.resultMap = myself.convertResultMap(resultItems)

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subject, content)
	}

	return nil
}

func (myself *baseFuture) convertResultMap(resultItems []*schema.ResultItem) map[string]*models.BaseResult {

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

func (myself *baseFuture) send(subject string, content []byte) error {
	nc, err := nats.Connect(myself.mqbrokerUrl)
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
func (myself *baseFuture) doReply(resMsg *schema.ResMsg) []*schema.ResultItem {

	return resMsg.Response

}

func (myself *baseFuture) prepareRequest(subjectChann string, reqmsg *models.RequestMsg, timeout time.Duration) {
	myself.reqMsg = reqmsg
	myself.timeout = timeout
	myself.subjectChann = subjectChann
}

/**
 * call and public request
 */
func (myself *baseFuture) publishRequest(subjectChann string, reqmsg *models.RequestMsg) error {
	myself.reqMsg = reqmsg
	myself.subjectChann = subjectChann

	reqMsg := new(schema.ReqMsg)
	reqMsg.Id = myself.reqMsg.Id

	paramByteData, err := msgpack.Marshal(myself.reqMsg.Params)
	reqMsg.Params = paramByteData
	if err != nil {
		return err
	}

	// --- contruct request message ---
	byteData, err := reqMsg.Marshal(nil)

	myself.send(myself.subjectChann, byteData)

	return nil
}

/**
 * receive response message
 */
func (myself *baseFuture) newFutureReturnResult(resmap map[string]*models.BaseResult) *baseFutureReturnResult {

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
