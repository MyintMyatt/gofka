package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	cmodels "dev/orion/gofka/client/models"
)

type GofkaServerRule struct {
	AllowedClientPrefixCommand []string
}

type GofkaServer struct {
	Clients map[*cmodels.Client]struct{}
	Rule    *GofkaServerRule
	mu      sync.RWMutex

	clientConnectedChann    chan *cmodels.Client
	clientDisconnectedChann chan *cmodels.Client
	clientSentMessageChann  chan *cmodels.Client
}

func NewGofkaServer() *GofkaServer {
	gs := &GofkaServer{
		Clients: make(map[*cmodels.Client]struct{}, 10),
		Rule: &GofkaServerRule{
			AllowedClientPrefixCommand: []string{"/name", "/lu", "/msg", "/q"},
		},

		clientConnectedChann:    make(chan *cmodels.Client, 100),
		clientDisconnectedChann: make(chan *cmodels.Client, 100),
		clientSentMessageChann:  make(chan *cmodels.Client, 100),
	}

	go gs.startRouterEvent()
	return gs
}

func (gs *GofkaServer) HandleConnection(conn net.Conn) {
	client := &cmodels.Client{
		Name:    "anon@" + conn.RemoteAddr().String(),
		Conn:    conn,
		MsgType: cmodels.ClientConnected,
	}

	gs.clientConnectedChann <- client

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", client.Name, err)
			gs.clientDisconnectedChann <- client
			return
		}

		msg := strings.TrimSpace(line)
		if msg == "" {
			continue
		}

		client.Msg = msg
		client.MsgType = cmodels.ClientSendMessage
		gs.clientSentMessageChann <- client
	}
}

func (gs *GofkaServer) startRouterEvent() {
	registerWorkers := 20
	for range registerWorkers {
		go gs.registerWorker()
	}

	broadcastWorkers := 20
	for range broadcastWorkers {
		go gs.broadCastWroker()
	}

	unRegisterWorkers := 20
	for range unRegisterWorkers {
		go gs.unRegisterWorker()
	}
}

// register
func (gs *GofkaServer) registerWorker() {
	for client := range gs.clientConnectedChann {
		gs.mu.Lock()
		gs.Clients[client] = struct{}{}
		gs.mu.Unlock()

		welcome := "##########\n" +
			"------Welcome to Gofka!------\n" +
			"------Announcement sign [***]------\n" +
			"------Warning sign [###]------\n" +
			"Avaliable Commands:\n" +
			"	/name <yourname> - set your display name\n" +
			"	/lu - show user lists\n" +
			"	/quit - leave the chat.\n" +
			"###########\n"
		client.Conn.Write([]byte(welcome))
		log.Printf("Client registered %s (%s)\n", client.Name, client.Conn.RemoteAddr())
	}
}

// unregister
func (gs *GofkaServer) unRegisterWorker() {
	for client := range gs.clientDisconnectedChann {
		gs.mu.Lock()

		if _, ok := gs.Clients[client]; ok {
			delete(gs.Clients, client)
			gs.mu.Unlock()

			_ = client.Conn.Close()
			log.Printf("Client unregistered: %s", client.Name)

			gs.broadCastSystemWide(fmt.Sprintf("*** Hey guys, unforturnely %s left the chat! ***\n", client.Name))
		} else {
			gs.mu.Unlock()
		}
	}
}

// broad cast
func (gs *GofkaServer) broadCastWroker() {
	for client := range gs.clientSentMessageChann {
		msg := client.Msg

		if strings.HasPrefix(msg, "/") {
			parts := strings.SplitN(msg, " ", 2)
			cmd := strings.ToLower(parts[0])

			switch cmd {
			case "/name":
				if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
					_, _ = client.Conn.Write([]byte("Usage: /name <yourname>\n"))
					continue
				}

				newName := strings.TrimSpace(parts[1])
				oldName := client.Name

				gs.mu.Lock()
				client.Name = newName
				gs.mu.Unlock()

				_, _ = client.Conn.Write([]byte(fmt.Sprintf("Name set to: %s\n", newName)))
				gs.broadCastSystemWide(fmt.Sprintf("*** %s is now known as %s ***\n", oldName, newName))
				continue
			case "/q":
				_, _ = client.Conn.Write([]byte("Goodbye!\n"))
				gs.clientDisconnectedChann <- client
				continue
			case "/lu":
				_, _ = client.Conn.Write([]byte(" --- Active User List ---\n"))
				gs.mu.Lock()
				count := 1
				for cl := range gs.Clients {
					client.Conn.Write([]byte(fmt.Sprintf("#[%d]: %s\n", count, cl.Name)))
					count++
				}
				gs.mu.Unlock()
				_, _ = client.Conn.Write([]byte("-------------------------\n"))
				continue
			default:
				_, _ = client.Conn.Write([]byte("### Unknow command. Available: /name, /q, /lu\n"))
			}
		}

		go gs.broadcast(client, msg)
	}
}

func (gs *GofkaServer) broadcast(sender *cmodels.Client, msg string) {
	formatted := fmt.Sprintf("[%s]: %s\n", sender.Name, msg)
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for c := range gs.Clients {
		if c == sender {
			continue
		}

		if _, err := c.Conn.Write([]byte(formatted)); err != nil {
			go func(cl *cmodels.Client) {
				gs.clientDisconnectedChann <- cl
			}(c)
		}
	}
}

func (gs *GofkaServer) broadCastSystemWide(msg string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for c := range gs.Clients {
		_, _ = c.Conn.Write([]byte(msg))
	}
}
