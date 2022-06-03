package crc

func Crc(bytes []byte) uint32 {
	var poly uint32 = (1 << 5) | (1 << 12)
	var bitOutMask uint32 = 1 << 16
	var state uint32 = ^uint32(0)
	for _, b := range bytes {
		// Loop for each bit
		for i := uint32(0); i < 8; i++ {
			state <<= 1
			if (state & bitOutMask) != 0 {
				state |= 1
			}
			if b&(1<<(7-i)) != 0 {
				state ^= 1
			}
			if state&1 != 0 {
				state ^= poly
			}
		}
	}
	return (^state) & 0xffff
}
