package jamulusprotocol

import (
	"bytes"
	"encoding/binary"
)

// ClientListItem
type ClientListItem struct {
	// 8-bit channel ID
	ChannelId uint8

	// 16-bit country ID
	ClientId uint16

	// 32-bit instrument ID
	InstrumentId uint32

	// 8-bit skill level
	SkillLevel uint8

	// Name
	Name string

	// City
	City string
}

func ParseClientListItem(buf *bytes.Buffer) (o *ClientListItem, err error) {
	o = &ClientListItem{}

	// Read channel ID
	err = binary.Read(buf, binary.LittleEndian, &o.ChannelId)
	if err != nil {
		return o, err
	}

	// Read country ID
	err = binary.Read(buf, binary.LittleEndian, &o.ClientId)
	if err != nil {
		return o, err
	}

	// Read instrument ID
	err = binary.Read(buf, binary.LittleEndian, &o.InstrumentId)
	if err != nil {
		return o, err
	}

	// Read skill level
	err = binary.Read(buf, binary.LittleEndian, &o.SkillLevel)
	if err != nil {
		return o, err
	}

	// Skip 4 zeroes
	buf.Next(4)

	// Read name length
	var nameLength uint16
	err = binary.Read(buf, binary.LittleEndian, &nameLength)
	if err != nil {
		return o, err
	}

	// Read name
	o.Name = string(buf.Next(int(nameLength)))

	// Read city length
	var cityLength uint16
	err = binary.Read(buf, binary.LittleEndian, &cityLength)
	if err != nil {
		return o, err
	}

	// Read city
	o.City = string(buf.Next(int(cityLength)))

	return
}
