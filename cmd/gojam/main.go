package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("go-jam!")

	pc, err := net.ListenPacket("udp", ":22199")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 8192)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		if buf[0] == 0x00 && buf[1] == 0x00 {
			// If first 2 bytes is 0x00 0x00, then it's a control message.
			// fmt.Printf("received %d bytes from %s\n", n, addr)
			reader := bytes.NewReader(buf[2:])

			// +-------------+------------+------------+------------------+ ...
			// | 2 bytes TAG | 2 bytes ID | 1 byte cnt | 2 bytes length n | ...
			// +-------------+------------+------------+------------------+ ...
			//     ... --------------+-------------+
			//     ...  n bytes data | 2 bytes CRC |
			//     ... --------------+-------------+
			var tag uint16
			var id uint16
			var cnt uint8
			var length uint16
			var data []byte
			var crc uint16
			err := binary.Read(reader, binary.LittleEndian, &tag)
			if err != nil {
				log.Fatal(err)
			}
			err = binary.Read(reader, binary.LittleEndian, &id)
			if err != nil {
				log.Fatal(err)
			}
			err = binary.Read(reader, binary.LittleEndian, &cnt)
			if err != nil {
				log.Fatal(err)
			}
			err = binary.Read(reader, binary.LittleEndian, &length)
			if err != nil {
				log.Fatal(err)
			}
			_, err = reader.Read(data)
			if err != nil {
				log.Fatal(err)
			}
			err = binary.Read(reader, binary.LittleEndian, &crc)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("received %d bytes from %s [tag: %d]\n", n, addr, tag)
		} else {
			// It is an audio packet.
			fmt.Printf("received %d bytes from %s\n", n, addr)
		}
	}
}
