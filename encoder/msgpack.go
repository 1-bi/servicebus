package encoder

import "github.com/1-bi/servicebus/schema"

const (
	ENCODER_TYPE_MSGPACK = 2
)

// MsgPackEncoder define encoder instance with messagepack
type MsgPackEncoder struct {
}

func (myself *MsgPackEncoder) GetType() byte {
	return ENCODER_TYPE_MSGPACK
}

func (myself *MsgPackEncoder) Encode(reqMsg *schema.ReqMsg) ([]byte, error) {
	return nil, nil
}

func (myself *MsgPackEncoder) Decode(inputContent []byte, resultObj interface{}) {

}
