package crc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrc(t *testing.T) {
	bytes := []byte{0x00, 0x00, 0x19, 0x00, 0x03, 0x1a, 0x00, 0xd3, 0x00, 0x07, 0x00, 0x00, 0x00, 0x01, 0x0f, 0x00, 0x64, 0x74, 0x69, 0x6e, 0x74, 0x68, 0x20, 0x2f, 0x20, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x00, 0x00}
	crc := Crc(bytes)
	assert.Equal(t, uint32(0x4D1C), crc)
}
