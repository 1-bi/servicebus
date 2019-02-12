package servicebus

// define encode and decode interface
type MessageEncoder interface {

	// --- encode object to byte ----
	Encode() ([]byte, error)

	// --- decode object from byte ----
	Decode(inputContent []byte, resultObj interface{})
}
