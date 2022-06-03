package audio

// #cgo CFLAGS: -I${SRCDIR}/../../libs
// #cgo LDFLAGS: -L${SRCDIR}/../../libs -lopus
// #include "opus/include/opus_custom.h"
import "C"
