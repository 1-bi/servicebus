package encoder

import "github.com/1-bi/servicebus/schema"

// GencodeEncoder define gencode encoder
type GencodeEncoder struct {
}

func (myself *GencodeEncoder) GetType() byte {
	return 0
}

func (myself *GencodeEncoder) Encode(reqMsg *schema.ReqMsg) ([]byte, error) {
	return nil, nil
}

func (myself *GencodeEncoder) Decode(inputContent []byte, resultObj interface{}) {

}
