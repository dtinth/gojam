package audio

// #cgo CFLAGS: -I${SRCDIR}/../../libs
// #cgo LDFLAGS: -L${SRCDIR}/../../libs/opus/.libs -lopus
// #include "opus/include/opus_custom.h"
import "C"
import "fmt"

type OpusDecoder struct {
	mode    *C.OpusCustomMode
	decoder *C.OpusCustomDecoder
}

func CreateDecoder() (*OpusDecoder, error) {
	var err C.int
	mode := C.opus_custom_mode_create(48000, 128, &err)
	if mode == nil {
		return nil, fmt.Errorf("opus_custom_mode_create: %d", err)
	}
	decoder := C.opus_custom_decoder_create(mode, 2, &err)
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
