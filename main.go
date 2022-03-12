package main

import (
	"fmt"
	"log"
)
import "go.bug.st/serial"

func main() {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
	}
	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write([]byte{0x7E, 0x01, 0x01, 0x00, 0xFE, 0x0D})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)

	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		for i := 0; i < n; i++ {
			fmt.Printf("0x%02X ", buff[i])
		}
		//fmt.Printf("0x%X ", string(buff[:n]))
	}

}
