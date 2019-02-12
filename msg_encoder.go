package servicebus

import "github.com/1-bi/servicebus/schema"

// define encode and decode interface
type MessageEncoder interface {

	// --- encode object to byte ----
	Encode(reqMsg *schema.ReqMsg) ([]byte, error)

	// --- decode object from byte ----
	Decode(inputContent []byte, resultObj interface{})
}
