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

type SkillLevelId uint8

const (
	SkillNone         SkillLevelId = 0
	SkillBeginner     SkillLevelId = 1
	SkillIntermediate SkillLevelId = 2
	SkillExpert       SkillLevelId = 3
)

type InstrumentId uint32

const (
	InstrumentNone             InstrumentId = 0
	InstrumentDrums            InstrumentId = 1
	InstrumentDjembe           InstrumentId = 2
	InstrumentElectricGuitar   InstrumentId = 3
	InstrumentAcousticGuitar   InstrumentId = 4
	InstrumentBassGuitar       InstrumentId = 5
	InstrumentKeyboard         InstrumentId = 6
	InstrumentSynthesizer      InstrumentId = 7
	InstrumentGrandPiano       InstrumentId = 8
	InstrumentAccordion        InstrumentId = 9
	InstrumentVocal            InstrumentId = 10
	InstrumentMicrophone       InstrumentId = 11
	InstrumentHarmonica        InstrumentId = 12
	InstrumentTrumpet          InstrumentId = 13
	InstrumentTrombone         InstrumentId = 14
	InstrumentFrenchHorn       InstrumentId = 15
	InstrumentTuba             InstrumentId = 16
	InstrumentSaxophone        InstrumentId = 17
	InstrumentClarinet         InstrumentId = 18
	InstrumentFlute            InstrumentId = 19
	InstrumentViolin           InstrumentId = 20
	InstrumentCello            InstrumentId = 21
	InstrumentDoubleBass       InstrumentId = 22
	InstrumentRecorder         InstrumentId = 23
	InstrumentStreamer         InstrumentId = 24
	InstrumentListener         InstrumentId = 25
	InstrumentGuitarVocal      InstrumentId = 26
	InstrumentKeyboardVocal    InstrumentId = 27
	InstrumentBodhran          InstrumentId = 28
	InstrumentBassoon          InstrumentId = 29
	InstrumentOboe             InstrumentId = 30
	InstrumentHarp             InstrumentId = 31
	InstrumentViola            InstrumentId = 32
	InstrumentCongas           InstrumentId = 33
	InstrumentBongo            InstrumentId = 34
	InstrumentVocalBass        InstrumentId = 35
	InstrumentVocalTenor       InstrumentId = 36
	InstrumentVocalAlto        InstrumentId = 37
	InstrumentVocalSoprano     InstrumentId = 38
	InstrumentBanjo            InstrumentId = 39
	InstrumentMandolin         InstrumentId = 40
	InstrumentUkulele          InstrumentId = 41
	InstrumentBassUkulele      InstrumentId = 42
	InstrumentVocalBaritone    InstrumentId = 43
	InstrumentVocalLead        InstrumentId = 44
	InstrumentMountainDulcimer InstrumentId = 45
	InstrumentScratching       InstrumentId = 46
	InstrumentRapping          InstrumentId = 47
	InstrumentVibraphone       InstrumentId = 48
	InstrumentConductor        InstrumentId = 49
)
