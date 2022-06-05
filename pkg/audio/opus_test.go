package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	decoder, err := CreateDecoder(2)
	if err != nil {
		t.Fatal(err)
	}
	var packet [165]byte
	packet[0] = 0x04
	packet[1] = 0xff
	packet[2] = 0xfe
	var output [960]int16
	output[0] = 0x01
	frames := decoder.Decode(packet[:], output[:])
	assert.Equal(t, 128, frames)
	assert.Equal(t, int16(0), output[0])
}
