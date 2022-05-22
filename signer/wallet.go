package wallet

import (
	"errors"
)

var (
	// ErrInputOutOfRange is returned by signing functions if the input index given is
	// less than zero or out of range of inputs in the transaction.
	ErrInputOutOfRange = errors.New("input index out of range for this transaction")
)

// type Wallet struct {
// 	key       []byte
// 	chainCode []byte
// }
