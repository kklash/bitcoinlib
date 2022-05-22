// Package signer provides high-level transaction signing helpers.
package signer

import (
	"errors"
)

// ErrInputOutOfRange is returned by signing functions if the input index given is
// less than zero or out of range of inputs in the transaction.
var ErrInputOutOfRange = errors.New("input index out of range for this transaction")
