package servicebus

// define encode and decode interface
type MessageEncoder interface {
	GetType() byte

	// --- encode object to byte ----
	//Encode(reqMsg *schema.ReqMsg) ([]byte, error)

	// --- decode object from byte ----
	Decode(inputContent []byte, resultObj interface{})
}
