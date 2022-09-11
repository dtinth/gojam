package gojam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dtinth/gojam/pkg/jamulusaudio"
	"github.com/dtinth/gojam/pkg/jamulusprotocol"
	"github.com/dtinth/gojam/pkg/jitterbuffer"
)

type Client struct {
	conn        net.Conn
	nextCounter uint8
	decoder     *jamulusaudio.OpusDecoder
	buffer      *jitterbuffer.JitterBuffer
	closed      bool
	info        jamulusprotocol.ChannelInfo

	// Handle PCM data. This is called when a new PCM data is available.
	// The PCM data is in stereo, 48kHz, 16-bit signed integer.
	HandlePCM func([]int16)

	// Print debug logging message
	DebugLog func(string)
}

// Creates a client
func NewClient(serverAddress string) (c *Client, err error) {
	c = &Client{}
	c.info = jamulusprotocol.ChannelInfo{
		Name:       "gojam",
		Country:    0,
		City:       "",
		Instrument: jamulusprotocol.InstrumentListener,
		SkillLevel: jamulusprotocol.SkillIntermediate,
	}

	// Create a decoder
	c.decoder, err = jamulusaudio.CreateDecoder(2)
	if err != nil {
		return nil, err
	}

	// Create a jitter buffer
	c.buffer = jitterbuffer.NewJitterBuffer(96)

	// Connect to the server
	conn, err := net.Dial("udp", serverAddress)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	// Print local address and remote address
	c.debug("Connected!")
	c.debug("Local address: %s", conn.LocalAddr())
	c.debug("Remote address: %s", conn.RemoteAddr())
	go c.sendSilence()
	go c.handleIncomingPackets()
	return
}

// Prints a debug logging message
func (c *Client) debug(format string, args ...interface{}) {
	if c.DebugLog != nil {
		c.DebugLog(fmt.Sprintf(format, args...))
	} else {
		format = "[gojam.Client DEBUG] " + format + "\n"
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

// Sends a silence packet to the server every 100ms
func (c *Client) sendSilence() {
	silence := jamulusaudio.NewSilentOpusStream()
	for !c.closed {
		packet := silence.Next()
		_, err := c.conn.Write(packet[:])
		if err != nil {
			c.debug("Error writing packet: %s", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Handles incoming packets from the server
func (c *Client) handleIncomingPackets() {
	buf := make([]byte, 8192)
	audioBuf := make([]int16, 8192)
	for !c.closed {
		n, err := c.conn.Read(buf)
		if err != nil {
			c.debug("Error reading from server: %s", err)
			continue
		}

		// If first 2 bytes is 0x00 0x00, then it's a protocol message.
		if buf[0] == 0x00 && buf[1] == 0x00 {
			message, err := jamulusprotocol.ParseMessage(buf)
			if err != nil {
				c.debug("Error parsing message: %s", err)
			} else {
				c.debug("Received message: %s", message)
				if message.Id != jamulusprotocol.Ackn && message.Id < 1000 {
					c.acknowledgeMessage(message.Id, message.Counter)
				}
				c.handleProtocolMessage(message)
			}
		} else {
			c.handleAudioPacket(buf[:n], audioBuf)
		}
	}
}

// Sends an acknowledgement message to the server
func (c *Client) acknowledgeMessage(id jamulusprotocol.MsgId, counter uint8) {
	// Generate a message data, which is just the ID
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(id))

	message := jamulusprotocol.Message{Id: jamulusprotocol.Ackn, Counter: counter, Data: data}
	c.sendMessage(message)
}

// Sends a message back to the server
func (c *Client) sendMessage(message jamulusprotocol.Message) {
	packet := message.Bytes()
	_, err := c.conn.Write(packet[:])
	if err != nil {
		c.debug("Error writing packet: %s", err)
	}
	c.debug("Sent message: %s", message)
}

// Handles a protocol message
func (c *Client) handleProtocolMessage(message jamulusprotocol.Message) {
	switch message.Id {
	case jamulusprotocol.ReqNetwTransportProps:
		c.sendNetwTransportProps()
	case jamulusprotocol.ReqJittBufSize:
		c.sendJittBufSize()
	case jamulusprotocol.ReqChannelInfos:
		c.sendChannelInfos()
	case jamulusprotocol.ConnClientsList:
		c.handleConnClientsList(message.Data)
	}
}

// Sends a NetwTransportProps message to the server
func (c *Client) sendNetwTransportProps() {
	// Create a buffer to hold the message data
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint32(166))   // Packet size
	binary.Write(&buf, binary.LittleEndian, uint16(2))     // Block size
	binary.Write(&buf, binary.LittleEndian, uint8(2))      // Channels
	binary.Write(&buf, binary.LittleEndian, uint32(48000)) // Sample rate
	binary.Write(&buf, binary.LittleEndian, uint16(2))     // Codec: Opus
	binary.Write(&buf, binary.LittleEndian, uint16(1))     // Flags: Add sequence number
	binary.Write(&buf, binary.LittleEndian, uint32(0))     // Other codec options

	// Assert that the buffer is 19 bytes long
	if buf.Len() != 19 {
		panic(fmt.Sprintf("Buffer length is %d, expected 19", buf.Len()))
	}

	message := jamulusprotocol.Message{
		Id:      jamulusprotocol.NetwTransportProps,
		Counter: c.nextCounterValue(),
		Data:    buf.Bytes(),
	}
	c.sendMessage(message)
}

// Sends a JittBufSize message to the server
func (c *Client) sendJittBufSize() {
	// Create a buffer to hold the message data
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, 4)
	message := jamulusprotocol.Message{
		Id:      jamulusprotocol.JittBufSize,
		Counter: c.nextCounterValue(),
		Data:    buf,
	}
	c.sendMessage(message)
}

// Sends a ChannelInfos message to the server
func (c *Client) sendChannelInfos() {
	message := jamulusprotocol.Message{
		Id:      jamulusprotocol.ChannelInfos,
		Counter: c.nextCounterValue(),
		Data:    c.info.Bytes(),
	}
	c.sendMessage(message)
}

// Update the client's channel info
func (c *Client) UpdateChannelInfo(info jamulusprotocol.ChannelInfo) {
	c.info = info
	c.sendChannelInfos()
}

// Handles a ConnClientsList message
func (c *Client) handleConnClientsList(data []byte) {
	// Create a buffer from the data
	buf := bytes.NewBuffer(data)

	// Read the clients until we reach the end of the buffer
	for buf.Len() > 0 {
		client, err := jamulusprotocol.ParseClientListItem(buf)
		if err != nil {
			c.debug("Error parsing client: %s", err)
			break
		}
		c.debug("Client list item: %s", client)
		c.unmute(client.ChannelId)
	}
}

// Unmutes a channel
func (c *Client) unmute(channelId uint8) {
	// Create a buffer to hold the message data
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint8(channelId)) // Channel ID
	binary.Write(&buf, binary.LittleEndian, uint16(0x8000))   // Gain

	message := jamulusprotocol.Message{
		Id:      jamulusprotocol.ChannelGain,
		Counter: c.nextCounterValue(),
		Data:    buf.Bytes(),
	}
	c.sendMessage(message)
}

// Get the next counter value
func (c *Client) nextCounterValue() uint8 {
	val := c.nextCounter
	c.nextCounter++
	return val
}

// Handles an audio packet
func (c *Client) handleAudioPacket(packet []byte, audioBuf []int16) {
	// Ensure that the packet is the correct size
	expectedSize := 332
	if len(packet) != expectedSize {
		c.debug("Audio packet is %d bytes, expected %d", len(packet), expectedSize)
		return
	}
	c.handleOpusPacket(packet[0:165], packet[165], audioBuf)
	c.handleOpusPacket(packet[166:331], packet[331], audioBuf)
}

// Handles an Opus packet
func (c *Client) handleOpusPacket(packet []byte, sequenceNum uint8, audioBuf []int16) {
	opusPacket := c.buffer.PutIn(packet, sequenceNum)
	if opusPacket != nil {
		samples := c.decoder.Decode(packet, audioBuf)
		if samples > 0 && c.HandlePCM != nil {
			c.HandlePCM(audioBuf[:samples*2])
		}
	}
}

// Close the client
func (c *Client) Close() {
	if c.closed {
		return
	}
	message := jamulusprotocol.Message{
		Id:      jamulusprotocol.ClmDisconnection,
		Counter: c.nextCounterValue(),
		Data:    []byte{},
	}
	c.sendMessage(message)
	c.closed = true
	c.conn.Close()
}
