package gojam

import (
	"fmt"
	"net"
	"time"

	"github.com/dtinth/gojam/pkg/jamulusaudio"
	"github.com/dtinth/gojam/pkg/jamulusprotocol"
)

type Client struct {
	conn net.Conn
}

// Creates a client
func NewClient(serverAddress string) (c *Client, err error) {
	c = &Client{}

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
	format = "DEBUG [gojamclient] " + format + "\n"
	fmt.Printf(format, args...)
}

// Sends a silence packet to the server every 100ms
func (c *Client) sendSilence() {
	silence := jamulusaudio.NewSilentOpusStream()
	for {
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
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			c.debug("Error reading from server: %s", err)
		}

		// If first 2 bytes is 0x00 0x00, then it's a control message.
		if buf[0] == 0x00 && buf[1] == 0x00 {
			msgId := jamulusprotocol.MsgId(int(buf[2]) + int(buf[3])*256)
			c.debug("Received control message: %s", msgId.String())
		} else {
			c.debug("Received audio packet of %d bytes", n)
		}
	}
}
