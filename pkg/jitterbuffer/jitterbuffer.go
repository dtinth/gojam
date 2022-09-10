// Package jitterbuffer implements a jitter buffer for audio frames
package jitterbuffer

// A Frame is a single audio frame
type Frame struct {
	SequenceNumber uint8
	Payload        []byte
}

// A JitterBuffer is a buffer for audio frames
type JitterBuffer struct {
	Size              int
	Frames            []Frame
	LatestSequenceNum uint8
}

// NewJitterBuffer creates a new jitter buffer
func NewJitterBuffer(size int) *JitterBuffer {
	return &JitterBuffer{
		Size:   size,
		Frames: make([]Frame, 0, size),
	}
}

// PutIn puts a frame into the jitter buffer, and returns the oldest frame if the buffer filled
func (j *JitterBuffer) PutIn(payload []byte, sequenceNum uint8) []byte {
	if len(j.Frames) == j.Size {
		// Pick the oldest frame and return it
		var oldestFrame *Frame
		for i := range j.Frames {
			if oldestFrame == nil || distance(j.LatestSequenceNum, j.Frames[i].SequenceNumber) > distance(j.LatestSequenceNum, oldestFrame.SequenceNumber) {
				oldestFrame = &j.Frames[i]
			}
		}

		// Retrive the payload
		oldestPayload := oldestFrame.Payload

		// Put a new frame in its place
		j.LatestSequenceNum = sequenceNum
		oldestFrame.SequenceNumber = sequenceNum
		oldestFrame.Payload = payload

		return oldestPayload
	}

	// Add the frame to the buffer
	j.Frames = append(j.Frames, Frame{
		SequenceNumber: sequenceNum,
		Payload:        payload,
	})
	return nil
}

func distance(latestSequenceNum uint8, sequenceNum uint8) int16 {
	diff := int16(latestSequenceNum) - int16(sequenceNum)
	if diff < -128 {
		diff += 256
	}
	if diff > 128 {
		diff -= 256
	}
	return diff
}
