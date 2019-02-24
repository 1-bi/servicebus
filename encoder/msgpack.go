package encoder

import "github.com/1-bi/servicebus/schema"

// MsgPackEncoder define encoder instance with messagepack
type MsgPackEncoder struct {
}

func (myself *MsgPackEncoder) GetType() byte {
	return 1
}

func (myself *MsgPackEncoder) Encode(reqMsg *schema.ReqMsg) ([]byte, error) {
	return nil, nil
}

func (myself *MsgPackEncoder) Decode(inputContent []byte, resultObj interface{}) {

}
