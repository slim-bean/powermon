package main

import (
	"encoding/json"
	"fmt"
	"log"
	"powermon/pkg/eg4"
	"time"
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

	buff := make([]byte, 500)
	for {

		n, err := port.Write([]byte{0x7E, 0x01, 0x01, 0x00, 0xFE, 0x0D})
		if err != nil {
			log.Fatal(err)
		}

		n, err = port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}

		if buff[0] == 0x7E && buff[n-1] == 0x0D {
			//for i := 0; i < n; i++ {
			//	fmt.Printf("%d:0x%02X ", i, buff[i])
			//}
			//fmt.Printf("\n\n")
			packet, err := eg4.Parse(buff)
			if err != nil {
				fmt.Printf("Error parsing packet: %v\n", err)
			}
			ps, err := json.Marshal(packet)
			if err != nil {
				fmt.Printf("Failed to marshal packet to json: %v\n", err)
			}
			fmt.Println(string(ps))
		} else {
			fmt.Printf("Did not receive valid packet from battery, ignoring this poll.\n")
		}

		time.Sleep(1000 * time.Millisecond)
		//fmt.Printf("0x%X ", string(buff[:n]))
	}

}
