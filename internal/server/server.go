package server

import (
	"log"
	"net"
	"sync"

	cmodels "dev/orion/gofka/client/models"
)

type GofkaServerRule struct {
	AllowedClientPrefixCommand []string
}

type GofkaServer struct {
	Clients map[*cmodels.Client]struct{}
	Rule    *GofkaServerRule
}

func NewGofkaServer() *GofkaServer {
	return &GofkaServer{
		Clients: make(map[*cmodels.Client]struct{}, 10),
		Rule: &GofkaServerRule{
			AllowedClientPrefixCommand: []string{"@username", "@msg"},
		},
	}
}

func (gs *GofkaServer) HandleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	clientChannel := make(chan *cmodels.Client, 200)

	go func() {
		clientChannel <- &cmodels.Client{
			Name:    "client: " + conn.RemoteAddr().String(),
			Conn:    conn,
			MsgType: cmodels.ClientConnected,
		}
	}()

	gs.handle(clientChannel)
}

func (gs *GofkaServer) handle(clientChannel chan *cmodels.Client) {

	newClientChann := make(chan *cmodels.Client, 50)
	deleteClientChann := make(chan *cmodels.Client, 50)
	broadcastMsgChann := make(chan *cmodels.Client, 100)

	client := <-clientChannel

	switch client.MsgType {
	case cmodels.ClientConnected:
		newClientChann <- client
	case cmodels.ClientDisConnected:
		deleteClientChann <- client
	case cmodels.ClientSendMessage:
		broadcastMsgChann <- client
	}

	var wg sync.WaitGroup
	numWorkers := 5
	for w := 1; w <= numWorkers; w++ {
		go gs.register()
	}
	log.Printf("Client %s connect to Gofka\n", client.Name)
}

func (gs *GofkaServer) register() {
}
