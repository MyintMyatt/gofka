package main

import (
	"encoding/binary"
	"fmt"
)

// "log"
// "net"
//
// b "dev/orion/gofka/internal/server"

func main() {
	value := uint32(0x78132456)
	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf, value)
	fmt.Printf("Big-endian:    % x\n", buf)

	binary.LittleEndian.PutUint32(buf, value)
	fmt.Printf("Little-endian: % x\n", buf)
	// log.Printf("Welcome from Gofka\n")
	// server, err := net.Listen("tcp", ":6969")
	// if err != nil {
	// 	log.Fatalf("Gofka Server can't start due to : %v\n", err)
	// }
	//
	// broker := b.NewBroker()
	//
	// for {
	// 	conn, err := server.Accept()
	// 	if err != nil {
	// 		log.Printf("Client connect error : %v\n", err)
	// 	}
	//
	// 	go broker.HandleConnection(conn)
	// }
}
