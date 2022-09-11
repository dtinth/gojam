package jamulusprotocol

import (
	"bytes"
	"encoding/binary"
)

type ChannelInfo struct {
	// Name
	Name string

	// Country
	Country uint16

	// City
	City string

	// Instrument
	Instrument InstrumentId

	// Skill Level
	SkillLevel SkillLevelId
}

func (c *ChannelInfo) Bytes() []byte {
	nameBytes := []byte(c.Name)
	cityBytes := []byte(c.City)

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(c.Country))      // Country
	binary.Write(&buf, binary.LittleEndian, uint32(c.Instrument))   // Listener
	binary.Write(&buf, binary.LittleEndian, uint8(c.SkillLevel))    // Skill Level
	binary.Write(&buf, binary.LittleEndian, uint16(len(nameBytes))) // Name length
	binary.Write(&buf, binary.LittleEndian, nameBytes)              // Name
	binary.Write(&buf, binary.LittleEndian, uint16(len(cityBytes))) // City length
	binary.Write(&buf, binary.LittleEndian, cityBytes)              // City
	return buf.Bytes()
}
