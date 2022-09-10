package jamulusprotocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dtinth/gojam/pkg/jamuluscrc"
)

// A Jamulus protocol message.
type Message struct {
	// 16-bit message ID
	Id MsgId

	// 8-bit message counter
	Counter uint8

	// Message data
	Data []byte
}

// Parses a Jamulus protocol message.
func ParseMessage(data []byte) (m Message, err error) {
	m = Message{}
	reader := bytes.NewReader(data)

	// First 2 bytes should be 0x00 0x00
	var tag uint16
	err = binary.Read(reader, binary.LittleEndian, &tag)
	if err != nil {
		return
	}
	if tag != 0x0000 {
		err = ErrInvalidTag
		return
	}

	// Read message ID
	var id uint16
	err = binary.Read(reader, binary.LittleEndian, &id)
	m.Id = MsgId(id)
	if err != nil {
		return
	}

	// Read message counter
	err = binary.Read(reader, binary.LittleEndian, &m.Counter)
	if err != nil {
		return
	}

	// Read message length
	var length uint16
	err = binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		return
	}

	// Read message data
	m.Data = make([]byte, length)
	_, err = reader.Read(m.Data)
	if err != nil {
		return
	}

	// Read message CRC
	var crc uint16
	err = binary.Read(reader, binary.LittleEndian, &crc)
	if err != nil {
		return
	}

	// Validate CRC
	expected := jamuluscrc.Crc(data[:length+7])
	if crc != expected {
		err = fmt.Errorf("crc mismatch: expected %d, got %d", crc, jamuluscrc.Crc(data[:len(data)-2]))
		return
	}

	return
}

var ErrInvalidTag = errors.New("invalid tag")

// Formats a Jamulus protocol message as a string.
func (m Message) String() string {
	return fmt.Sprintf("Message(Id: %d [%s], Counter: %d)", m.Id, m.Id.String(), m.Counter)
}

// Serializes a Jamulus protocol message.
func (m Message) ToPacket() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(0x0000)) // Tag
	binary.Write(&buf, binary.LittleEndian, uint16(m.Id))   // ID
	binary.Write(&buf, binary.LittleEndian, m.Counter)      // Counter
	binary.Write(&buf, binary.LittleEndian, uint16(len(m.Data)))
	buf.Write(m.Data)
	crc := jamuluscrc.Crc(buf.Bytes())
	binary.Write(&buf, binary.LittleEndian, crc)
	return buf.Bytes()
}
