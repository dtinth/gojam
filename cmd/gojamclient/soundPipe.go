package main

import (
	"bytes"
	"encoding/binary"
)

type soundPipe struct {
	vadEnabled     bool
	hp             int
	HandlePCMBytes func([]byte)
}

func newSoundPipe(vadEnabled bool) *soundPipe {
	return &soundPipe{
		vadEnabled: vadEnabled,
	}
}

func (p *soundPipe) WritePCM(pcm []int16) {
	var buf bytes.Buffer
	activityDetected := false
	for _, sample := range pcm {
		binary.Write(&buf, binary.LittleEndian, sample)
		if sample != 0 {
			activityDetected = true
		}
	}
	if p.vadEnabled {
		if activityDetected {
			p.hp = 100
		} else {
			p.hp--
			if p.hp < 0 {
				p.hp = 0
			}
			if p.hp == 0 {
				return
			}
		}
	}
	if p.HandlePCMBytes != nil {
		p.HandlePCMBytes(buf.Bytes())
	}
}
