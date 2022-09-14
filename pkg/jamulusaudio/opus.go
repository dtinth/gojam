package jamulusaudio

// #cgo LDFLAGS: -lopus -lm
// #include "opus_custom.h"
// #include "opus_defines.h"
// int gojam_opus_disable_vbr(OpusCustomEncoder* encoder) { return opus_custom_encoder_ctl(encoder, OPUS_SET_VBR(0)); }
// int gojam_opus_set_low_delay(OpusCustomEncoder* encoder) { return opus_custom_encoder_ctl(encoder, OPUS_SET_APPLICATION(OPUS_APPLICATION_RESTRICTED_LOWDELAY)); }
// int gojam_opus_set_low_complexity(OpusCustomEncoder* encoder) { return opus_custom_encoder_ctl(encoder, OPUS_SET_COMPLEXITY(1)); }
import "C"
import "fmt"

// An Opus decoder for the Jamulus protocol.
type OpusDecoder struct {
	mode    *C.OpusCustomMode
	decoder *C.OpusCustomDecoder
}

// Creates an Opus decoder.
func CreateDecoder(numChannels int) (*OpusDecoder, error) {
	var err C.int
	mode := C.opus_custom_mode_create(48000, 128, &err)
	if mode == nil {
		return nil, fmt.Errorf("opus_custom_mode_create: %d", err)
	}
	decoder := C.opus_custom_decoder_create(mode, C.int(numChannels), &err)
	if decoder == nil {
		return nil, fmt.Errorf("opus_custom_decoder_create: %d", err)
	}
	return &OpusDecoder{mode, decoder}, nil
}

// Decodes an Opus packet into PCM samples.
func (d *OpusDecoder) Decode(packet []byte, output []int16) int {
	result := C.opus_custom_decode(d.decoder, (*C.uchar)(&packet[0]), C.int(len(packet)), (*C.opus_int16)(&output[0]), C.int(len(output)/2))
	return int(result)
}

// Destroys an Opus decoder.
func (d *OpusDecoder) Destroy() {
	C.opus_custom_decoder_destroy(d.decoder)
	C.opus_custom_mode_destroy(d.mode)
}

// An Opus encode for the Jamulus protocol.
type OpusEncoder struct {
	mode    *C.OpusCustomMode
	encoder *C.OpusCustomEncoder
}

// Creates an Opus encoder.
func CreateEncoder(numChannels int) (*OpusEncoder, error) {
	var err C.int
	mode := C.opus_custom_mode_create(48000, 128, &err)
	if mode == nil {
		return nil, fmt.Errorf("opus_custom_mode_create: %d", err)
	}
	encoder := C.opus_custom_encoder_create(mode, C.int(numChannels), &err)
	if encoder == nil {
		return nil, fmt.Errorf("opus_custom_encoder_create: %d", err)
	}

	C.gojam_opus_disable_vbr(encoder)
	C.gojam_opus_set_low_delay(encoder)
	C.gojam_opus_set_low_complexity(encoder)

	return &OpusEncoder{mode, encoder}, nil
}
