package server

import (
	"net"
	"sync"

	log "dev/orion/gofka/internal/log"
	proto "dev/orion/gofka/internal/protocol"
)

// /
// / [4 bytes total length][2 bytes API key][2 bytes version][4 bytes correlation ID][payload...]
// /
type GofkaBroker struct {
	logs map[string]*log.PartitionLog
	mu   sync.RWMutex
}

func NewBroker() *GofkaBroker {
	return &GofkaBroker{}
}

// /
// / client handling
// /
func (b *GofkaBroker) HandleConnection(conn *net.Conn) {

}

// /
// / Data processing
// /
func (b *GofkaBroker) HandleProduce(req *proto.ProduceRequest) *proto.ProduceResponse {
	return &proto.ProduceResponse{}
}

func (b *GofkaBroker) HandleFetch(req *proto.FetchRequest) *proto.FetchResponse {
	return &proto.FetchResponse{}
}
