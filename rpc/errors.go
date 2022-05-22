package rpc

import (
	"errors"
	"fmt"
)

type ErrRPCFailure struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *ErrRPCFailure) Error() string {
	subError, ok := ErrorsByCode[err.Code]
	if ok {
		return fmt.Sprintf("%s - %s", subError, err.Message)
	}

	return err.Message
}

func (err *ErrRPCFailure) Unwrap() (subError error) {
	subError, _ = ErrorsByCode[err.Code]
	return
}

// ErrorsByCode is a map of error codes to errors. Messages and
// codes are taken from bitcoin/src/rpc/protocol.h
var ErrorsByCode = map[int]error{
	// General application defined errors
	-1:  errors.New("std::exception thrown in command handling"),
	-3:  errors.New("Unexpected type was passed as parameter"),
	-5:  errors.New("Invalid address or key"),
	-7:  errors.New("Ran out of memory during operation"),
	-8:  errors.New("Invalid, missing or duplicate parameter"),
	-20: errors.New("Database error"),
	-22: errors.New("Error parsing or validating structure in raw format"),
	-25: errors.New("General error during transaction or block submission"),
	-26: errors.New("Transaction or block was rejected by network rules"),
	-27: errors.New("Transaction already in chain"),
	-28: errors.New("Client still warming up"),
	-32: errors.New("RPC method is deprecated"),

	// P2P client errors
	-9:  errors.New("Bitcoin is not connected"),
	-10: errors.New("Still downloading initial blocks"),
	-23: errors.New("Node is already added"),
	-24: errors.New("Node has not been added before"),
	-29: errors.New("Node to disconnect not found in connected nodes"),
	-30: errors.New("Invalid IP/Subnet"),
	-31: errors.New("No valid connection manager instance found"),

	// Wallet errors
	-4:  errors.New("Unspecified problem with wallet (key not found etc.)"),
	-6:  errors.New("Not enough funds in wallet or account"),
	-11: errors.New("Invalid label name"),
	-12: errors.New("Keypool ran out, call keypoolrefill first"),
	-13: errors.New("Enter the wallet passphrase with walletpassphrase first"),
	-14: errors.New("The wallet passphrase entered was incorrect"),
	-15: errors.New("Command given in wrong wallet encryption state (encrypting an encrypted wallet etc.)"),
	-16: errors.New("Failed to encrypt the wallet"),
	-17: errors.New("Wallet is already unlocked"),
	-18: errors.New("Invalid wallet specified"),
	-19: errors.New("No wallet specified (error when there are multiple wallets loaded)"),
}
