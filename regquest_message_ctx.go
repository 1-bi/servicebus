package servicebus

type embeddedReqMsgContext struct {
	rawMsgBody []byte
}

func newEmbeddedReqMsgContext() *embeddedReqMsgContext {
	var ctx = new(embeddedReqMsgContext)
	return ctx
}

func (myself *embeddedReqMsgContext) setMsgRawBody(body []byte) {
	myself.rawMsgBody = body
}

func (myself *embeddedReqMsgContext) GetMsgRawBody() []byte {
	return myself.rawMsgBody
}

func (myself *embeddedReqMsgContext) GetResResult() Result {
	return nil
}
