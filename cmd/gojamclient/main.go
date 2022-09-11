package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dtinth/gojam/pkg/gojam"
)

func main() {
	// Parse command line arguments
	server := flag.String("server", "localhost:22124", "server to connect to")
	pcmout := flag.String("pcmout", "", "server to pipe PCM data to")
	name := flag.String("name", "", "musician name")

	flag.Parse()

	fmt.Println("gojam jamulus client")
	fmt.Println("server:", *server)

	client, err := gojam.NewClient(*server)
	if err != nil {
		panic(err)
	}
	fmt.Println("client:", client)

	info := client.GetChannelInfo()
	if *name != "" {
		info.Name = *name
		client.UpdateChannelInfo(info)
	}

	if *pcmout != "" {
		installPCMOut(client, *pcmout)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("exiting")
	client.Close()
}

func installPCMOut(client *gojam.Client, pcmout string) {
	conn, err := net.Dial("tcp", pcmout)
	if err != nil {
		panic(err)
	}
	outChan := make(chan []byte, 100)
	go func() {
		for {
			data := <-outChan
			_, err := conn.Write(data)
			if err != nil {
				client.Close()
				panic(err)
			}
		}
	}()
	client.HandlePCM = func(pcm []int16) {
		var buf bytes.Buffer
		for _, sample := range pcm {
			binary.Write(&buf, binary.LittleEndian, sample)
		}
		outChan <- buf.Bytes()
	}
}
