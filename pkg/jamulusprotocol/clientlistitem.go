package jamulusprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// ClientListItem
type ClientListItem struct {
	// 8-bit channel ID
	ChannelId uint8

	// 16-bit country ID
	CountryId uint16

	// 32-bit instrument ID
	InstrumentId InstrumentId

	// 8-bit skill level
	SkillLevel SkillLevelId

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
	err = binary.Read(buf, binary.LittleEndian, &o.CountryId)
	if err != nil {
		return o, err
	}

	// Read instrument ID
	var instrument uint32
	err = binary.Read(buf, binary.LittleEndian, &instrument)
	o.InstrumentId = InstrumentId(instrument)
	if err != nil {
		return o, err
	}

	// Read skill level
	var skillLevel uint8
	err = binary.Read(buf, binary.LittleEndian, &skillLevel)
	o.SkillLevel = SkillLevelId(skillLevel)
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

// Return a string representation of this object.
func (o *ClientListItem) String() string {
	return fmt.Sprintf(
		"ClientListItem{#%d \"%s\" from \"%s\", \"%d\" %s (%d), %s}",
		o.ChannelId,
		o.Name,
		o.City,
		o.CountryId,
		o.InstrumentId,
		o.InstrumentId,
		o.SkillLevel,
	)
}
