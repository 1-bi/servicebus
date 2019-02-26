package encoder

import "github.com/1-bi/servicebus/schema"

// GencodeEncoder define gencode encoder
const (
	ENCODER_TYPE_GENCODE = 1
)

type GencodeEncoder struct {
}

func (myself *GencodeEncoder) GetType() byte {
	return ENCODER_TYPE_GENCODE
}

func (myself *GencodeEncoder) Encode(reqMsg *schema.ReqMsg) ([]byte, error) {
	return nil, nil
}

func (myself *GencodeEncoder) Decode(inputContent []byte, resultObj interface{}) {

}
