package models

import (
	"net"
)

type MessageType int

const (
	ClientConnected MessageType = iota
	ClientDisConnected
	ClientSendMessage
)

type Client struct {
	Name    string
	Conn    net.Conn
	MsgType MessageType
}
