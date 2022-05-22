package blockscan

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/kklash/bitcoinlib/rpc"
)

func createScanner() (*BlockScanner, error) {
	username, password, err := rpc.FindCookie()
	if err != nil {
		if errors.Is(err, rpc.ErrCookieNotFound) {
			err = fmt.Errorf("Ensure bitcoind is running for this test to pass")
		}
		return nil, err
	}

	conn, err := rpc.NewConnection("http://127.0.0.1:8332", username, password)
	if err != nil {
		return nil, err
	}

	return NewBlockScanner(conn), nil
}

func TestBlockScanner_GetBlockByHeight(t *testing.T) {
	t.Skip()

	scanner, err := createScanner()
	if err != nil {
		t.Errorf("failed to create scanner: %s", err)
		return
	}

	block, err := scanner.GetBlockByHeight(586)
	if err != nil {
		t.Errorf("failed to get block: %s", err)
		return
	}

	expectedBlock := strings.Join(
		[]string{
			"AQAAADi6vJWGpfzWBxNXNJT0N358QBwzqiRymk9s/0YAAAAATVlpwNENzOYIaP7k1N6Aul7zirru",
			"2KddqmPkjJY9exlQR29J//8AHS2XkTcDAQAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
			"AAAAAP////8IBP//AB0CXQb/////AQDyBSoBAAAAQ0EEENrwSe9ALeC2rbqLD3w5K8+aY4URbvyL",
			"QUO4t6eEHn3nO0eP/hO2DFDqAeJLS0jCT14PvF1shDPHynw+06uBc6wAAAAAAQAAAAUPQPXmXhFe",
			"tL2zAH8PuL6qQEz3rkXeFgdOisybabvwwwAAAABIRzBEAiAQktpAr23qiry+77hYYzWybTnTa+m2",
			"w41snMGPIN1YhgIgRZZN55qQCPaNU/ybxY+eMLIkobmNv9pce3uGDzLGrvEB/////xu4dbJHMy5V",
			"hzHCxRD2EdPd6ZHqn+aTZb9EWgzNUTsZAAAAAElIMEUCIQCwodCgAlHFaAmlq117psvmi4LJv0+A",
			"buOcVorlN1cshAIgeBzmkBfsOy1vlv//TRnIDCJPQMc7jCbLpLMOf0FxV5sB/////yCZ4aktlMNf",
			"BkVoMlfEwlUWU4Xz6RKahf7Vo/PYZ8m2AAAAAElIMEUCIQDI6YD0PGFiMuLVnc4Ipe24SqoJFepJ",
			"eAqK82czAhYISgIgPMJijxb5lceq9hBMumSXGWOk4ITk+9C2vPgltHoJ+OMB/////1+3cMTecArK",
			"f3T15ilfJI7a+pQj5EbXb0ZQ35uQ+TmnAAAAAElIMEUCIHRajZnFH5j1yTuNL18UofLYzEL/cylk",
			"VoG8r+hGy/UNAiEAsk4xGGEp865syKIm0e2jiTc2UqnPIJVjH8xDRQZ8H/MB/////5aNTAlu6GEw",
			"eTXSHXl6kCtkfclw08g3TME1Ufg5ervYAAAAAElIMEUCIQDKZbPykHJNbFb8MzVw+jQvJHfzSyps",
			"k8Li1yFtn+kIjgIgd+JZop7R+Yj6srnyzhekpWogwYjK3HK8qU4GpzgmlmUB/////wEAuh3SBQAA",
			"AENBBJcwTv06sU0Ny/HpAQRaJfS1269XbQdFBv2N7UEium9r7A7UaYzg55KMDq+d37U4eSm11pfo",
			"Lnqr6+BMEOXIcWSsAAAAAAEAAAABDSa6V/+C/vy0OCa0UBkEPitu+aqBGLf3QxZ1hKf5yucAAAAA",
			"SUgwRQIgJP1zRd8rK9Dm+EFlKQRrfVK9pf/bcBRrxtcrG6c8q80CIQD/mcAwBsyPKNkuaG8K5kDS",
			"A5UXfzKdCp29Vg/SpVru5wH/////AQDyBSoBAAAAQ0EEiI2JDhvYTJ4qw2Opd0QUoIHrgFzSwNUu",
			"Se/HFw6/NC8c2yhKLi63VPyN1FJf4Mqj06UlIU0LUE3XU3ay9jgEqKwAAAAA",
		},
		"",
	)

	actualBlock := base64.StdEncoding.EncodeToString(block.Bytes())

	if actualBlock != expectedBlock {
		t.Errorf("fetched block does not match expected\nWanted %s\nGot    %s", expectedBlock, actualBlock)
		return
	}

	expectedHash := "000000000d0d23516c5efd3af4eb951603bb30b2c93884b522a318b30e918ee7"

	blockHash, err := block.Header.Hash()
	if err != nil {
		t.Errorf("failed to hash block header: %s", err)
		return
	}

	actualHash := hex.EncodeToString(blockHash[:])

	if actualHash != expectedHash {
		t.Errorf("block hash does not match\nWanted %s\nGot    %s", expectedHash, actualHash)
		return
	}
}
