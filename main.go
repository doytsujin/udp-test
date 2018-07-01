package main

import (
	"net"
	"fmt"
	"encoding/hex"
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
		fmt.Println(n, a, err )
		fmt.Println(hex.Dump(buf))
		data := binary.BigEndian.Uint32(buf[:4])
		fmt.Printf("%b %X %X \n", data, data, data|0xFF000000) //buf[1]&8, buf[1]&4
		//fmt.Printf("%s \n", buf)
		break
		//file.Write(buf)
	}


}
