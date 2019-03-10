package validation

type MsgCandidateVali struct {
	Header []byte `validate:"check-encoder-match"`
}
