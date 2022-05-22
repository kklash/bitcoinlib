package script

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kklash/bitcoinlib/constants"
)

func TestPushReadData(t *testing.T) {
	type Fixture struct {
		dataSize   int
		firstBytes []byte
	}

	dataFixtures := []Fixture{
		Fixture{0x00, []byte{0x00}},
		Fixture{0x01, []byte{0x01}},
		Fixture{0x10, []byte{0x10}},
		Fixture{0x51, []byte{constants.OP_PUSHDATA1, 0x51}},
		Fixture{0x4c, []byte{constants.OP_PUSHDATA1, 0x4c}},
		Fixture{0xff, []byte{constants.OP_PUSHDATA1, 0xff}},
		Fixture{0x0208, []byte{constants.OP_PUSHDATA2, 0x08, 0x02}},
	}

	for _, fixture := range dataFixtures {
		data := make([]byte, fixture.dataSize) // filled with zero bytes
		pushScript := PushData(data)
		if !bytes.HasPrefix(pushScript, fixture.firstBytes) {
			t.Errorf("push encoding failed for script length %d - wanted %x, got %x", fixture.dataSize, fixture.firstBytes, pushScript[:5])
		} else if !bytes.HasSuffix(pushScript, data) {
			t.Errorf("push encoding failed for script length %d - data was lost", fixture.dataSize)
		} else if len(pushScript) != len(fixture.firstBytes)+fixture.dataSize {
			t.Errorf("unexpected data length: wanted %d, got %d", len(fixture.firstBytes)+fixture.dataSize, len(pushScript))
		}

		readData, err := ReadData(bytes.NewReader(pushScript))
		if err != nil {
			t.Errorf("failed to decode data: %s", err)
		}

		if fmt.Sprintf("%x", data) != fmt.Sprintf("%x", readData) {
			t.Errorf("decoded data did not match\nwanted %s\ngot %s", fmt.Sprintf("%x", data), fmt.Sprintf("%x", readData))
		}
	}
}

func TestPushReadNumber(t *testing.T) {
	type Fixture struct {
		n     int64
		bytes []byte
	}

	fixtures := []Fixture{
		Fixture{1, []byte{constants.OP_1}},
		Fixture{0, []byte{constants.OP_0}},
		Fixture{12, []byte{constants.OP_12}},
		Fixture{16, []byte{constants.OP_16}},
		Fixture{17, []byte{0x01, 17}},                                                              // In little endian:
		Fixture{0x7f, []byte{0x01, 0x7f}},                                                          // 11111110
		Fixture{-0x7f, []byte{0x01, 0xff}},                                                         // 11111111 most significant bit set for negative
		Fixture{0xff, []byte{0x02, 0xff, 0x00}},                                                    // 11111111 00000000 extra zero byte indicates non-negative
		Fixture{-0xff, []byte{0x02, 0xff, 0x80}},                                                   // 11111111 00000001 most significant byte set
		Fixture{0x7fff, []byte{0x02, 0xff, 0x7f}},                                                  // 11111111 11111110
		Fixture{-0x7fff, []byte{0x02, 0xff, 0xff}},                                                 // 11111111 11111111
		Fixture{0xffff, []byte{0x03, 0xff, 0xff, 0x00}},                                            // 11111111 11111111 00000000
		Fixture{-0xffff, []byte{0x03, 0xff, 0xff, 0x80}},                                           // 11111111 11111111 00000001
		Fixture{0x7fffffffffffffff, []byte{0x08, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}},  // 11111111 ... 11111110
		Fixture{-0x7fffffffffffffff, []byte{0x08, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}}, // 11111111 ... 11111111
	}

	for _, fixture := range fixtures {
		data := PushNumber(fixture.n)
		if !bytes.Equal(data, fixture.bytes) {
			t.Errorf("push number encoding failed - wanted %x, got %x", fixture.bytes, data)
			continue
		}

		n, err := ReadNumber(bytes.NewReader(data))
		if err != nil {
			t.Errorf("failed to read number from encoded data: %s", err)
			continue
		}

		if n != fixture.n {
			t.Errorf("decoded number does not match - wanted %d , got %d", fixture.n, n)
		}
	}
}
