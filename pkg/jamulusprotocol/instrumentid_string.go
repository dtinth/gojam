// Code generated by "stringer -type=InstrumentId"; DO NOT EDIT.

package jamulusprotocol

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InstrumentNone-0]
	_ = x[InstrumentDrums-1]
	_ = x[InstrumentDjembe-2]
	_ = x[InstrumentElectricGuitar-3]
	_ = x[InstrumentAcousticGuitar-4]
	_ = x[InstrumentBassGuitar-5]
	_ = x[InstrumentKeyboard-6]
	_ = x[InstrumentSynthesizer-7]
	_ = x[InstrumentGrandPiano-8]
	_ = x[InstrumentAccordion-9]
	_ = x[InstrumentVocal-10]
	_ = x[InstrumentMicrophone-11]
	_ = x[InstrumentHarmonica-12]
	_ = x[InstrumentTrumpet-13]
	_ = x[InstrumentTrombone-14]
	_ = x[InstrumentFrenchHorn-15]
	_ = x[InstrumentTuba-16]
	_ = x[InstrumentSaxophone-17]
	_ = x[InstrumentClarinet-18]
	_ = x[InstrumentFlute-19]
	_ = x[InstrumentViolin-20]
	_ = x[InstrumentCello-21]
	_ = x[InstrumentDoubleBass-22]
	_ = x[InstrumentRecorder-23]
	_ = x[InstrumentStreamer-24]
	_ = x[InstrumentListener-25]
	_ = x[InstrumentGuitarVocal-26]
	_ = x[InstrumentKeyboardVocal-27]
	_ = x[InstrumentBodhran-28]
	_ = x[InstrumentBassoon-29]
	_ = x[InstrumentOboe-30]
	_ = x[InstrumentHarp-31]
	_ = x[InstrumentViola-32]
	_ = x[InstrumentCongas-33]
	_ = x[InstrumentBongo-34]
	_ = x[InstrumentVocalBass-35]
	_ = x[InstrumentVocalTenor-36]
	_ = x[InstrumentVocalAlto-37]
	_ = x[InstrumentVocalSoprano-38]
	_ = x[InstrumentBanjo-39]
	_ = x[InstrumentMandolin-40]
	_ = x[InstrumentUkulele-41]
	_ = x[InstrumentBassUkulele-42]
	_ = x[InstrumentVocalBaritone-43]
	_ = x[InstrumentVocalLead-44]
	_ = x[InstrumentMountainDulcimer-45]
	_ = x[InstrumentScratching-46]
	_ = x[InstrumentRapping-47]
	_ = x[InstrumentVibraphone-48]
	_ = x[InstrumentConductor-49]
}

const _InstrumentId_name = "InstrumentNoneInstrumentDrumsInstrumentDjembeInstrumentElectricGuitarInstrumentAcousticGuitarInstrumentBassGuitarInstrumentKeyboardInstrumentSynthesizerInstrumentGrandPianoInstrumentAccordionInstrumentVocalInstrumentMicrophoneInstrumentHarmonicaInstrumentTrumpetInstrumentTromboneInstrumentFrenchHornInstrumentTubaInstrumentSaxophoneInstrumentClarinetInstrumentFluteInstrumentViolinInstrumentCelloInstrumentDoubleBassInstrumentRecorderInstrumentStreamerInstrumentListenerInstrumentGuitarVocalInstrumentKeyboardVocalInstrumentBodhranInstrumentBassoonInstrumentOboeInstrumentHarpInstrumentViolaInstrumentCongasInstrumentBongoInstrumentVocalBassInstrumentVocalTenorInstrumentVocalAltoInstrumentVocalSopranoInstrumentBanjoInstrumentMandolinInstrumentUkuleleInstrumentBassUkuleleInstrumentVocalBaritoneInstrumentVocalLeadInstrumentMountainDulcimerInstrumentScratchingInstrumentRappingInstrumentVibraphoneInstrumentConductor"

var _InstrumentId_index = [...]uint16{0, 14, 29, 45, 69, 93, 113, 131, 152, 172, 191, 206, 226, 245, 262, 280, 300, 314, 333, 351, 366, 382, 397, 417, 435, 453, 471, 492, 515, 532, 549, 563, 577, 592, 608, 623, 642, 662, 681, 703, 718, 736, 753, 774, 797, 816, 842, 862, 879, 899, 918}

func (i InstrumentId) String() string {
	if i >= InstrumentId(len(_InstrumentId_index)-1) {
		return "InstrumentId(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _InstrumentId_name[_InstrumentId_index[i]:_InstrumentId_index[i+1]]
}