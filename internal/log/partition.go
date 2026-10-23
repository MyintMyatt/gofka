package log

import (
	"sync"

	protocol "dev/orion/gofka/internal/protocol"
)

type PartitionLog struct {
	topic      string
	partiton   int32
	dir        string
	mu         sync.Mutex
	nextOffset int64
}

func (p *PartitionLog) Append(msgs []protocol.Message) (baseOffset int64, err error) {
	return 1, nil
}

// func (p *PartitionLog) Read() ()
