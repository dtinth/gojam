package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type mp3Broadcaster struct {
	mutex       sync.Mutex
	connections map[string]chan []byte
	stdin       io.WriteCloser
}

func newMp3Broadcaster() *mp3Broadcaster {
	cmd := exec.Command("ffmpeg", "-f", "s16le", "-ar", "48000", "-ac", "2", "-i", "-", "-f", "mp3", "-b:a", "192k", "-")
	in, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	b := &mp3Broadcaster{
		connections: make(map[string]chan []byte),
		stdin:       in,
	}
	go func() {
		for {
			buf := make([]byte, 8192)
			n, err := out.Read(buf)
			if err != nil {
				panic(err)
			}
			if n > 0 {
				b.sendMp3Data(buf[:n])
			}
		}
	}()
	return b
}

func (b *mp3Broadcaster) sendMp3Data(out []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for id, ch := range b.connections {
		// Do not block
		select {
		case ch <- out:
		default:
			fmt.Fprintf(os.Stderr, "[mp3Broadcaster] drop %s\n", id)
		}
	}
}

func (b *mp3Broadcaster) Register(id string, out chan []byte) func() {
	fmt.Fprintf(os.Stderr, "[mp3Broadcaster] lock %s\n", id)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	fmt.Fprintf(os.Stderr, "[mp3Broadcaster] register %s\n", id)
	b.connections[id] = out
	return func() {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		fmt.Fprintf(os.Stderr, "[mp3Broadcaster] unregister %s\n", id)
		delete(b.connections, id)
	}
}

func (b *mp3Broadcaster) WritePCMBytes(pcmBytes []byte) {
	_, err := b.stdin.Write(pcmBytes)
	if err != nil {
		panic(err)
	}
}

func (b *mp3Broadcaster) GetListenerCount() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return len(b.connections)
}
