package address

import (
	"errors"
)

var (
	// ErrInvalidPublicKeyLength indicates an improper length of public key was passed
	// to an address-making function.
	ErrInvalidPublicKeyLength = errors.New("invalid public key length")

	// ErrInvalidHash indicates an incorrect length hash was passed to MakeFromHash.
	ErrInvalidHash = errors.New("incorrect hash length passed when making address from hash")

	// ErrInvalidAddressFormat indicates that an invalid address format string
	// was passed when creating a new address.
	ErrInvalidAddressFormat = errors.New("invalid address type passed; use constants declared in this package")

	// ErrNoSegwitSupport indicates the caller tried to make a segwit address
	// for a network which has no segwit support.
	ErrNoSegwitSupport = errors.New("cannot make segwit addresses for network with no segwit support")

	// ErrInvalidAddress indicates the caller passed an invalid address to an address decoding function.
	ErrInvalidAddress = errors.New("invalid address, cannot decode")
)
