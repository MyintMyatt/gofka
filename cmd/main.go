package main

import (
	"log"
	"net"

	gofka "dev/orion/gofka/internal/server"
)

func main() {

	log.Printf("Welcome from Gofka\n")
	server, err := net.Listen("tcp", ":6969")
	if err != nil {
		log.Fatalf("Gofka Server can't start due to : %v\n", err)
	}

	gofkaServer := gofka.NewGofkaServer()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Printf("Client connect error : %v\n", err)
		}

		go gofkaServer.HandleConnection(conn)
	}
}
