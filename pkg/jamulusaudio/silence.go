package jamulusaudio

// Generates a stream of silent Opus packets.
type SilentOpusStream struct {
	counter uint8
}

// Creates a new silent Opus stream.
func NewSilentOpusStream() *SilentOpusStream {
	return &SilentOpusStream{0}
}

// Generates the next packet in the stream.
func (s *SilentOpusStream) Next() [332]byte {
	var packet [332]byte
	write := func(slice []byte) {
		slice[0] = 0x04
		slice[1] = 0xff
		slice[2] = 0xfe
		s.counter = s.counter + 1
		slice[165] = s.counter
	}
	write(packet[:])
	write(packet[166:])
	return packet
}
