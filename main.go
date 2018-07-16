package main

import (
	"net"
	"fmt"
	"encoding/binary"
)

type Header struct {
	Value [4]byte
	SyncByte byte

}

func main(){
	// 239.255.10.160:5500
	addr, err := net.ResolveUDPAddr("udp", "239.255.10.160:5500")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	//file, err := os.Open("video.mp4")
	if err != nil {
		panic(err)
	}
	for {
		conn.SetReadBuffer(4048)
		buf := make([]byte,4048)
		n, a, err := conn.ReadFrom(buf)
		fmt.Printf("Get %d bytes packet from %s. Error %v \n", n, a, err )
		//fmt.Println(hex.Dump(buf))
		data := binary.BigEndian.Uint32(buf[:4])
		fmt.Printf("%b %X %X \n", data, data, data|0xFF000000) //buf[1]&8, buf[1]&4
		fmt.Printf("Sync byte: ASCII %c %t \n", buf[0], buf[0]==0x47)
		fmt.Printf("Transport Error Indicator: %b \n", data & 0x800000)
		fmt.Printf("Payload Unit Start Indicator : %b \n", data & 0x400000)
		fmt.Printf("Transport Priority: %b \n", data & 0x200000)
		fmt.Printf("PID: %b \n", data & 0x1fff00)
		fmt.Printf("Transport Scrambling Control : %b\n", data & 0xc0)
		switch data & 0xc0 {
		case 0x0:
			fmt.Println("Not scrambled")
		case 0x40:
			fmt.Println("Reserved for future use")
		case 0x80:
			fmt.Println("Scrambled with even key")
		case 0xC0:
			fmt.Println("Scrambled with odd key")
		}
		fmt.Printf("Adaptation field control : %b\n", data & 0x30)
		switch data & 0x30 {
		case 1<<4: //10000
			fmt.Println("no adaptation field, payload only")
		case 1<<5: //100000
			fmt.Println("adaptation field only, no payload")
		case 1<<5 | 1: // 110000
			fmt.Println("adaptation field followed by payload")
		default:
			fmt.Println("RESERVED for future use")
		}
		fmt.Printf("Continuity counter : %b\n", data & 0xf)
		// first 3 bytes of payload has to 0x000001 ; buf[4:8] is called the 32 bit start code.
		if buf[4] == 0x00 && buf[5] == 0x00 && buf[6] == 0x01 {
			fmt.Println("Packet start code prefix")
		}
		fmt.Printf("Stream id %x\n", buf[7])
		if buf[7] >= 0xE0 && 0xEF <= buf[7] {
			fmt.Println("Video stream")
		}
		if buf[7] >= 0xC0 && 0xDF <= buf[7] {
			fmt.Println("Audio stream")
		}
		lenPesPacket := binary.BigEndian.Uint16(buf[8:10])
		fmt.Printf("PES Packet length %d \n", lenPesPacket)
		//fmt.Printf("%s \n", buf)
		break
		//file.Write(buf)
	}


}
