package script

import (
	"fmt"

	"github.com/kklash/bitcoinlib/constants"
)

// ErrOpReturnDataTooLarge is returned by MakeOpReturn when the size of an OP_RETURN
// data payload exceeds constants.OpReturnMaxSize.
var ErrOpReturnDataTooLarge = fmt.Errorf("OP_RETURN payloads must be smaller than %d bytes", constants.OpReturnMaxSize)

// MakeOpReturn creates an OP_RETURN output script containing the given payload. Returns
// ErrOpReturnDataTooLarge if the payload is larger than constants.OpReturnMaxSize.
func MakeOpReturn(payload []byte) ([]byte, error) {
	if len(payload) > constants.OpReturnMaxSize {
		return nil, ErrOpReturnDataTooLarge
	}
	output := append([]byte{constants.OP_RETURN}, PushData(payload)...)
	return output, nil
}
