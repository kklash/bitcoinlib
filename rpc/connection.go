package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	// ErrInvalidCredentials is returned if an RPC request returns an authorization error.
	ErrInvalidCredentials = errors.New("invalid username+password for rpc.Connection")

	// ErrInvalidResponseFormat is returned when an RPC response cannot be parsed.
	ErrInvalidResponseFormat = errors.New("RPC returned incorrect response format")

	// ErrInvalidUrl is returned when initiating a connection to an unsupported URL.
	ErrInvalidUrl = errors.New("invalid URL protocol for Connection")
)

// Connection represents an HTTP connection to a Bitcoin node.
type Connection struct {
	URL                *url.URL
	HttpClient         *http.Client
	username, password string
	requestID          int
	requestIDMutex     *sync.Mutex
}

// NewConnection initiates a new Connection to the given URL.
func NewConnection(uri, username, password string) (*Connection, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return nil, ErrInvalidUrl
	}

	conn := &Connection{
		URL:            parsedURL,
		HttpClient:     http.DefaultClient,
		username:       username,
		password:       password,
		requestIDMutex: new(sync.Mutex),
	}

	return conn, nil
}

type rpcResponseMessage struct {
	Error  *ErrRPCFailure `json:"error"`
	Result any            `json:"result"`
}

// Request initiates a new RPC method call over the Connection and returns the
// RPC response's "result" property as a JSON-decoded any.
func (conn *Connection) Request(method string, params ...any) (any, error) {
	var result any
	if err := conn.RequestSetResult(&result, method, params...); err != nil {
		return nil, err
	}

	return result, nil
}

// RequestSetResult is similar to conn.Request, except the caller can
// pass a pointer into which the RPC response result is decoded.
func (conn *Connection) RequestSetResult(resultPtr any, method string, params ...any) error {
	conn.requestIDMutex.Lock()
	requestID := conn.requestID
	conn.requestID += 1
	conn.requestIDMutex.Unlock()

	reqBodyMap := map[string]any{
		"jsonrpc": "1.0",
		"id":      requestID,
		"method":  method,
		"params":  params,
	}

	reqBody := new(bytes.Buffer)

	if err := json.NewEncoder(reqBody).Encode(reqBodyMap); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, conn.URL.String(), reqBody)
	if err != nil {
		return err
	}

	req.SetBasicAuth(conn.username, conn.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return ErrInvalidCredentials
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	bodyString := string(bodyBytes)

	if bodyString == "Work queue depth exceeded" {
		fmt.Fprintln(
			os.Stderr,
			"bitcoin/rpc WARNING: Exceeding work queue depth for bitcoind.",
			"Request will be retried, but you should decrease parallel request",
			"load, or increase -rpcworkqueue and -rpcthreads settings in",
			"bitcoin.conf (defaults are 16 and 4 respectively).",
		)
		time.Sleep(100 * time.Millisecond)
		return conn.RequestSetResult(resultPtr, method, params...)
	}

	responseObj := &rpcResponseMessage{Result: resultPtr}

	if err := json.Unmarshal(bodyBytes, &responseObj); err != nil {
		if resp.StatusCode != 200 {
			return fmt.Errorf("%w: %s", ErrInvalidResponseFormat, bodyString)
		}

		return fmt.Errorf("%w: %s", ErrInvalidResponseFormat, err)
	}

	if responseObj.Error != nil {
		return responseObj.Error
	}

	if responseObj.Result != nil {
		return nil
	}

	return ErrInvalidResponseFormat
}
