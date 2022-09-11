package main

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/dtinth/gojam/pkg/gojam"
)

func main() {
	jamulusServerAddress := flag.String("jamulusserver", "localhost:22124", "Jamulus server")
	token := flag.String("token", "", "Bot Token")
	guild := flag.String("guild", "", "Guild")
	channel := flag.String("channel", "", "Channel")
	flag.Parse()

	if *token == "" {
		panic("No token provided")
	}
	if *guild == "" {
		panic("No guild provided")
	}
	if *channel == "" {
		panic("No channel provided")
	}

	flag.Parse()
	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer s.Close()

	j, err := gojam.NewClient(*jamulusServerAddress)
	if err != nil {
		panic(err)
	}
	defer j.Close()

	// We only really care about receiving voice state updates.
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildVoiceStates)

	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	v, err := s.ChannelVoiceJoin(*guild, *channel, false, true)
	if err != nil {
		fmt.Println("failed to join voice channel:", err)
		return
	}

	v.Speaking(true)

	audioChannel := make(chan []int16)
	var pcmBuffer PcmBuffer
	pcmBuffer.OnFull = func(buffer []int16) {
		audioChannel <- buffer
	}
	j.HandlePCM = func(pcm []int16) {
		pcmBuffer.Add(pcm)
	}

	dgvoice.SendPCM(v, audioChannel)
}

// Buffers 960*2 samples of audio
type PcmBuffer struct {
	buffer []int16
	index  int
	OnFull func(buffer []int16)
}

// Adds samples to the buffer, and when the buffer is filled, calls the OnFull callback.
func (b *PcmBuffer) Add(samples []int16) {
	if b.buffer == nil {
		b.buffer = make([]int16, 960*2)
	}
	for _, sample := range samples {
		b.buffer[b.index] = sample
		b.index++
		if b.index == len(b.buffer) {
			b.OnFull(b.buffer)
			b.index = 0
			b.buffer = make([]int16, 960*2)
		}
	}
}
